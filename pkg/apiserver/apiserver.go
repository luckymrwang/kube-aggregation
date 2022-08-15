/*
Copyright 2019 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apiserver

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	rt "runtime"
	"strconv"
	"sync"
	"time"

	"github.com/emicklei/go-restful"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	urlruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	unionauth "k8s.io/apiserver/pkg/authentication/request/union"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog"
	clusterv1alpha1 "kubesphere.io/api/cluster/v1alpha1"
	iamv1alpha2 "kubesphere.io/api/iam/v1alpha2"
	typesv1beta1 "kubesphere.io/api/types/v1beta1"
	runtimecache "sigs.k8s.io/controller-runtime/pkg/cache"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"kube-aggregation/pkg/apiserver/authentication/authenticators/basic"
	"kube-aggregation/pkg/apiserver/authentication/authenticators/jwt"
	"kube-aggregation/pkg/apiserver/authentication/request/anonymous"
	"kube-aggregation/pkg/apiserver/authentication/request/basictoken"
	"kube-aggregation/pkg/apiserver/authentication/request/bearertoken"
	"kube-aggregation/pkg/apiserver/authentication/token"
	"kube-aggregation/pkg/apiserver/authorization"
	"kube-aggregation/pkg/apiserver/authorization/authorizer"
	"kube-aggregation/pkg/apiserver/authorization/authorizerfactory"
	"kube-aggregation/pkg/apiserver/authorization/path"
	"kube-aggregation/pkg/apiserver/authorization/rbac"
	unionauthorizer "kube-aggregation/pkg/apiserver/authorization/union"
	apiserverconfig "kube-aggregation/pkg/apiserver/config"
	"kube-aggregation/pkg/apiserver/dispatch"
	"kube-aggregation/pkg/apiserver/filters"
	"kube-aggregation/pkg/apiserver/request"
	"kube-aggregation/pkg/informers"
	clusterkapisv1alpha1 "kube-aggregation/pkg/kapis/cluster/v1alpha1"
	configv1alpha2 "kube-aggregation/pkg/kapis/config/v1alpha2"
	"kube-aggregation/pkg/kapis/crd"
	iamapi "kube-aggregation/pkg/kapis/iam/v1alpha2"
	meteringv1alpha1 "kube-aggregation/pkg/kapis/metering/v1alpha1"
	monitoringv1alpha3 "kube-aggregation/pkg/kapis/monitoring/v1alpha3"
	"kube-aggregation/pkg/kapis/oauth"
	resourcesv1alpha2 "kube-aggregation/pkg/kapis/resources/v1alpha2"
	resourcev1alpha3 "kube-aggregation/pkg/kapis/resources/v1alpha3"
	"kube-aggregation/pkg/kapis/version"
	"kube-aggregation/pkg/models/auth"
	"kube-aggregation/pkg/models/iam/am"
	"kube-aggregation/pkg/models/iam/group"
	"kube-aggregation/pkg/models/iam/im"
	"kube-aggregation/pkg/models/resources/v1alpha3/loginrecord"
	"kube-aggregation/pkg/models/resources/v1alpha3/user"
	"kube-aggregation/pkg/simple/client/cache"
	"kube-aggregation/pkg/simple/client/k8s"
	"kube-aggregation/pkg/simple/client/monitoring"
	"kube-aggregation/pkg/utils/clusterclient"
	"kube-aggregation/pkg/utils/iputil"
	"kube-aggregation/pkg/utils/metrics"
)

var initMetrics sync.Once

type APIServer struct {
	// number of kubesphere apiserver
	ServerCount int

	Server *http.Server

	Config *apiserverconfig.Config

	// webservice container, where all webservice defines
	container *restful.Container

	// kubeClient is a collection of all kubernetes(include CRDs) objects clientset
	KubernetesClient k8s.Client

	// informerFactory is a collection of all kubernetes(include CRDs) objects informers,
	// mainly for fast query
	InformerFactory informers.InformerFactory

	// cache is used for short lived objects, like session
	CacheClient cache.Interface

	// monitoring client set
	MonitoringClient monitoring.Interface

	MetricsClient monitoring.Interface

	// controller-runtime cache
	RuntimeCache runtimecache.Cache

	// entity that issues tokens
	Issuer token.Issuer

	// controller-runtime client
	RuntimeClient runtimeclient.Client

	ClusterClient clusterclient.ClusterClients
}

func (s *APIServer) PrepareRun(stopCh <-chan struct{}) error {
	s.container = restful.NewContainer()
	s.container.Filter(logRequestAndResponse)
	s.container.Router(restful.CurlyRouter{})
	s.container.RecoverHandler(func(panicReason interface{}, httpWriter http.ResponseWriter) {
		logStackOnRecover(panicReason, httpWriter)
	})

	s.installKubeSphereAPIs(stopCh)
	s.installCRDAPIs()
	s.installMetricsAPI()
	s.container.Filter(monitorRequest)

	for _, ws := range s.container.RegisteredWebServices() {
		klog.V(2).Infof("%s", ws.RootPath())
	}

	s.Server.Handler = s.container

	s.buildHandlerChain(stopCh)

	return nil
}

func monitorRequest(r *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	start := time.Now()
	chain.ProcessFilter(r, response)
	reqInfo, exists := request.RequestInfoFrom(r.Request.Context())
	if exists && reqInfo.APIGroup != "" {
		RequestCounter.WithLabelValues(reqInfo.Verb, reqInfo.APIGroup, reqInfo.APIVersion, reqInfo.Resource, strconv.Itoa(response.StatusCode())).Inc()
		elapsedSeconds := time.Now().Sub(start).Seconds()
		RequestLatencies.WithLabelValues(reqInfo.Verb, reqInfo.APIGroup, reqInfo.APIVersion, reqInfo.Resource).Observe(elapsedSeconds)
	}
}

func (s *APIServer) installMetricsAPI() {
	initMetrics.Do(registerMetrics)
	metrics.Defaults.Install(s.container)
}

// Install all kubesphere api groups
// Installation happens before all informers start to cache objects, so
//   any attempt to list objects using listers will get empty results.
func (s *APIServer) installKubeSphereAPIs(stopCh <-chan struct{}) {
	imOperator := im.NewOperator(s.KubernetesClient.KubeSphere(),
		user.New(s.InformerFactory.KubeSphereSharedInformerFactory(),
			s.InformerFactory.KubernetesSharedInformerFactory()),
		loginrecord.New(s.InformerFactory.KubeSphereSharedInformerFactory()),
		s.Config.AuthenticationOptions)
	amOperator := am.NewOperator(s.KubernetesClient.KubeSphere(),
		s.KubernetesClient.Kubernetes(),
		s.InformerFactory)
	rbacAuthorizer := rbac.NewRBACAuthorizer(amOperator)

	urlruntime.Must(configv1alpha2.AddToContainer(s.container, s.Config))
	urlruntime.Must(resourcev1alpha3.AddToContainer(s.container, s.InformerFactory, s.RuntimeCache))
	urlruntime.Must(monitoringv1alpha3.AddToContainer(s.container, s.KubernetesClient.Kubernetes(), s.MonitoringClient, s.MetricsClient, s.InformerFactory, s.RuntimeClient))
	urlruntime.Must(meteringv1alpha1.AddToContainer(s.container, s.KubernetesClient.Kubernetes(), s.MonitoringClient, s.InformerFactory, s.RuntimeCache, s.Config.MeteringOptions, s.RuntimeClient))
	urlruntime.Must(resourcesv1alpha2.AddToContainer(s.container, s.KubernetesClient.Kubernetes(), s.InformerFactory,
		s.KubernetesClient.Master()))
	urlruntime.Must(clusterkapisv1alpha1.AddToContainer(s.container,
		s.KubernetesClient.KubeSphere(),
		s.InformerFactory.KubernetesSharedInformerFactory(),
		s.InformerFactory.KubeSphereSharedInformerFactory(),
		s.Config.MultiClusterOptions.ProxyPublishService,
		s.Config.MultiClusterOptions.ProxyPublishAddress,
		s.Config.MultiClusterOptions.AgentImage))
	urlruntime.Must(iamapi.AddToContainer(s.container, imOperator, amOperator,
		group.New(s.InformerFactory, s.KubernetesClient.KubeSphere(), s.KubernetesClient.Kubernetes()),
		rbacAuthorizer))

	userLister := s.InformerFactory.KubeSphereSharedInformerFactory().Iam().V1alpha2().Users().Lister()
	urlruntime.Must(oauth.AddToContainer(s.container, imOperator,
		auth.NewTokenOperator(s.CacheClient, s.Issuer, s.Config.AuthenticationOptions),
		auth.NewPasswordAuthenticator(s.KubernetesClient.KubeSphere(), userLister, s.Config.AuthenticationOptions),
		auth.NewOAuthAuthenticator(s.KubernetesClient.KubeSphere(), userLister, s.Config.AuthenticationOptions),
		auth.NewLoginRecorder(s.KubernetesClient.KubeSphere(), userLister),
		s.Config.AuthenticationOptions))
	urlruntime.Must(version.AddToContainer(s.container, s.KubernetesClient.Kubernetes().Discovery()))
}

// installCRDAPIs Install CRDs to the KAPIs with List and Get options
func (s *APIServer) installCRDAPIs() {
	crds := &extv1.CustomResourceDefinitionList{}
	// TODO Maybe we need a better label name
	urlruntime.Must(s.RuntimeClient.List(context.TODO(), crds, runtimeclient.MatchingLabels{"kubesphere.io/resource-served": "true"}))
	urlruntime.Must(crd.AddToContainer(s.container, s.RuntimeClient, s.RuntimeCache, crds))
}

func (s *APIServer) Run(ctx context.Context) (err error) {

	err = s.waitForResourceSync(ctx)
	if err != nil {
		return err
	}

	shutdownCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-ctx.Done()
		_ = s.Server.Shutdown(shutdownCtx)
	}()

	klog.V(0).Infof("Start listening on %s", s.Server.Addr)
	if s.Server.TLSConfig != nil {
		err = s.Server.ListenAndServeTLS("", "")
	} else {
		err = s.Server.ListenAndServe()
	}

	return err
}

func (s *APIServer) buildHandlerChain(stopCh <-chan struct{}) {
	requestInfoResolver := &request.RequestInfoFactory{
		APIPrefixes:          sets.NewString("api", "apis", "kapis", "kapi"),
		GrouplessAPIPrefixes: sets.NewString("api", "kapi"),
		GlobalResources: []schema.GroupResource{
			iamv1alpha2.Resource(iamv1alpha2.ResourcesPluralUser),
			iamv1alpha2.Resource(iamv1alpha2.ResourcesPluralGlobalRole),
			iamv1alpha2.Resource(iamv1alpha2.ResourcesPluralGlobalRoleBinding),
			clusterv1alpha1.Resource(clusterv1alpha1.ResourcesPluralCluster),
			resourcev1alpha3.Resource(clusterv1alpha1.ResourcesPluralCluster),
		},
	}

	handler := s.Server.Handler
	handler = filters.WithKubeAPIServer(handler, s.KubernetesClient.Config(), &errorResponder{})

	var authorizers authorizer.Authorizer

	switch s.Config.AuthorizationOptions.Mode {
	case authorization.AlwaysAllow:
		authorizers = authorizerfactory.NewAlwaysAllowAuthorizer()
	case authorization.AlwaysDeny:
		authorizers = authorizerfactory.NewAlwaysDenyAuthorizer()
	default:
		fallthrough
	case authorization.RBAC:
		excludedPaths := []string{"/oauth/*", "/kapis/config.kubesphere.io/*", "/kapis/version", "/kapis/metrics"}
		pathAuthorizer, _ := path.NewAuthorizer(excludedPaths)
		amOperator := am.NewReadOnlyOperator(s.InformerFactory)
		authorizers = unionauthorizer.New(pathAuthorizer, rbac.NewRBACAuthorizer(amOperator))
	}

	handler = filters.WithAuthorization(handler, authorizers)
	if s.Config.MultiClusterOptions.Enable {
		clusterDispatcher := dispatch.NewClusterDispatch(s.ClusterClient)
		handler = filters.WithMultipleClusterDispatcher(handler, clusterDispatcher)
	}

	userLister := s.InformerFactory.KubeSphereSharedInformerFactory().Iam().V1alpha2().Users().Lister()
	loginRecorder := auth.NewLoginRecorder(s.KubernetesClient.KubeSphere(), userLister)

	// authenticators are unordered
	authn := unionauth.New(anonymous.NewAuthenticator(),
		basictoken.New(basic.NewBasicAuthenticator(auth.NewPasswordAuthenticator(
			s.KubernetesClient.KubeSphere(),
			userLister,
			s.Config.AuthenticationOptions),
			loginRecorder)),
		bearertoken.New(jwt.NewTokenAuthenticator(
			auth.NewTokenOperator(s.CacheClient, s.Issuer, s.Config.AuthenticationOptions),
			userLister)))
	handler = filters.WithAuthentication(handler, authn)
	handler = filters.WithRequestInfo(handler, requestInfoResolver)

	s.Server.Handler = handler
}

func isResourceExists(apiResources []v1.APIResource, resource schema.GroupVersionResource) bool {
	for _, apiResource := range apiResources {
		if apiResource.Name == resource.Resource {
			return true
		}
	}
	return false
}

type informerForResourceFunc func(resource schema.GroupVersionResource) (interface{}, error)

func waitForCacheSync(discoveryClient discovery.DiscoveryInterface, sharedInformerFactory informers.GenericInformerFactory, informerForResourceFunc informerForResourceFunc, GVRs map[schema.GroupVersion][]string, stopCh <-chan struct{}) error {
	for groupVersion, resourceNames := range GVRs {
		var apiResourceList *v1.APIResourceList
		var err error
		err = retry.OnError(retry.DefaultRetry, func(err error) bool {
			return !errors.IsNotFound(err)
		}, func() error {
			apiResourceList, err = discoveryClient.ServerResourcesForGroupVersion(groupVersion.String())
			return err
		})
		if err != nil {
			return fmt.Errorf("failed to fetch group version resources %s: %s", groupVersion, err)
		}
		for _, resourceName := range resourceNames {
			groupVersionResource := groupVersion.WithResource(resourceName)
			if !isResourceExists(apiResourceList.APIResources, groupVersionResource) {
				klog.Warningf("resource %s not exists in the cluster", groupVersionResource)
			} else {
				// reflect.ValueOf(sharedInformerFactory).MethodByName("ForResource").Call([]reflect.Value{reflect.ValueOf(groupVersionResource)})
				if _, err = informerForResourceFunc(groupVersionResource); err != nil {
					return fmt.Errorf("failed to create informer for %s: %s", groupVersionResource, err)
				}
			}
		}
	}
	sharedInformerFactory.Start(stopCh)
	sharedInformerFactory.WaitForCacheSync(stopCh)
	return nil
}

func (s *APIServer) waitForResourceSync(ctx context.Context) error {
	klog.V(0).Info("Start cache objects")

	stopCh := ctx.Done()
	// resources we have to create informer first
	k8sGVRs := map[schema.GroupVersion][]string{
		{Group: "", Version: "v1"}: {
			"namespaces",
			"nodes",
			"resourcequotas",
			"pods",
			"services",
			"persistentvolumeclaims",
			"persistentvolumes",
			"secrets",
			"configmaps",
			"serviceaccounts",
		},
		{Group: "rbac.authorization.k8s.io", Version: "v1"}: {
			"roles",
			"rolebindings",
			"clusterroles",
			"clusterrolebindings",
		},
		{Group: "apps", Version: "v1"}: {
			"deployments",
			"daemonsets",
			"replicasets",
			"statefulsets",
			"controllerrevisions",
		},
		{Group: "storage.k8s.io", Version: "v1"}: {
			"storageclasses",
		},
		{Group: "batch", Version: "v1"}: {
			"jobs",
		},
		{Group: "batch", Version: "v1beta1"}: {
			"cronjobs",
		},
		{Group: "networking.k8s.io", Version: "v1"}: {
			"ingresses",
			"networkpolicies",
		},
		{Group: "autoscaling", Version: "v2beta2"}: {
			"horizontalpodautoscalers",
		},
	}

	if err := waitForCacheSync(s.KubernetesClient.Kubernetes().Discovery(),
		s.InformerFactory.KubernetesSharedInformerFactory(),
		func(resource schema.GroupVersionResource) (interface{}, error) {
			return s.InformerFactory.KubernetesSharedInformerFactory().ForResource(resource)
		},
		k8sGVRs, stopCh); err != nil {
		return err
	}

	ksGVRs := map[schema.GroupVersion][]string{
		{Group: "iam.kubesphere.io", Version: "v1alpha2"}: {
			"users",
			"globalroles",
			"globalrolebindings",
			"groups",
			"groupbindings",
			"workspaceroles",
			"workspacerolebindings",
			"loginrecords",
		},
		{Group: "cluster.kubesphere.io", Version: "v1alpha1"}: {
			"clusters",
		},
	}

	// federated resources on cached in multi cluster setup
	if s.Config.MultiClusterOptions.Enable {
		ksGVRs[typesv1beta1.SchemeGroupVersion] = []string{
			typesv1beta1.ResourcePluralFederatedClusterRole,
			typesv1beta1.ResourcePluralFederatedClusterRoleBindingBinding,
			typesv1beta1.ResourcePluralFederatedNamespace,
			typesv1beta1.ResourcePluralFederatedService,
			typesv1beta1.ResourcePluralFederatedDeployment,
			typesv1beta1.ResourcePluralFederatedSecret,
			typesv1beta1.ResourcePluralFederatedConfigmap,
			typesv1beta1.ResourcePluralFederatedStatefulSet,
			typesv1beta1.ResourcePluralFederatedIngress,
			typesv1beta1.ResourcePluralFederatedPersistentVolumeClaim,
			typesv1beta1.ResourcePluralFederatedApplication,
		}
	}

	if err := waitForCacheSync(s.KubernetesClient.Kubernetes().Discovery(),
		s.InformerFactory.KubeSphereSharedInformerFactory(),
		func(resource schema.GroupVersionResource) (interface{}, error) {
			return s.InformerFactory.KubeSphereSharedInformerFactory().ForResource(resource)
		},
		ksGVRs, stopCh); err != nil {
		return err
	}

	go s.RuntimeCache.Start(ctx)
	s.RuntimeCache.WaitForCacheSync(ctx)

	klog.V(0).Info("Finished caching objects")
	return nil

}

func logStackOnRecover(panicReason interface{}, w http.ResponseWriter) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("recover from panic situation: - %v\r\n", panicReason))
	for i := 2; ; i += 1 {
		_, file, line, ok := rt.Caller(i)
		if !ok {
			break
		}
		buffer.WriteString(fmt.Sprintf("    %s:%d\r\n", file, line))
	}
	klog.Errorln(buffer.String())

	headers := http.Header{}
	if ct := w.Header().Get("Content-Type"); len(ct) > 0 {
		headers.Set("Accept", ct)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal server error"))
}

func logRequestAndResponse(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	start := time.Now()
	chain.ProcessFilter(req, resp)

	// Always log error response
	logWithVerbose := klog.V(4)
	if resp.StatusCode() > http.StatusBadRequest {
		logWithVerbose = klog.V(0)
	}

	logWithVerbose.Infof("%s - \"%s %s %s\" %d %d %dms",
		iputil.RemoteIp(req.Request),
		req.Request.Method,
		req.Request.URL,
		req.Request.Proto,
		resp.StatusCode(),
		resp.ContentLength(),
		time.Since(start)/time.Millisecond,
	)
}

type errorResponder struct{}

func (e *errorResponder) Error(w http.ResponseWriter, req *http.Request, err error) {
	klog.Error(err)
	responsewriters.InternalError(w, req, err)
}
