/*
Copyright 2020 The KubeSphere Authors.

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

package resource

import (
	"errors"
	"kube-aggregation/pkg/models/resources/v1alpha3/persistentvolume"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	monitoringdashboardv1alpha2 "kubesphere.io/monitoring-dashboard/api/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	"kube-aggregation/pkg/api"
	"kube-aggregation/pkg/apiserver/query"
	"kube-aggregation/pkg/informers"
	"kube-aggregation/pkg/models/resources/v1alpha3"
	"kube-aggregation/pkg/models/resources/v1alpha3/application"
	"kube-aggregation/pkg/models/resources/v1alpha3/clusterdashboard"
	"kube-aggregation/pkg/models/resources/v1alpha3/clusterrole"
	"kube-aggregation/pkg/models/resources/v1alpha3/clusterrolebinding"
	"kube-aggregation/pkg/models/resources/v1alpha3/configmap"
	"kube-aggregation/pkg/models/resources/v1alpha3/daemonset"
	"kube-aggregation/pkg/models/resources/v1alpha3/dashboard"
	"kube-aggregation/pkg/models/resources/v1alpha3/deployment"
	"kube-aggregation/pkg/models/resources/v1alpha3/globalrole"
	"kube-aggregation/pkg/models/resources/v1alpha3/globalrolebinding"
	"kube-aggregation/pkg/models/resources/v1alpha3/group"
	"kube-aggregation/pkg/models/resources/v1alpha3/groupbinding"
	"kube-aggregation/pkg/models/resources/v1alpha3/ingress"
	"kube-aggregation/pkg/models/resources/v1alpha3/job"
	"kube-aggregation/pkg/models/resources/v1alpha3/loginrecord"
	"kube-aggregation/pkg/models/resources/v1alpha3/namespace"
	"kube-aggregation/pkg/models/resources/v1alpha3/networkpolicy"
	"kube-aggregation/pkg/models/resources/v1alpha3/node"
	"kube-aggregation/pkg/models/resources/v1alpha3/pod"
	"kube-aggregation/pkg/models/resources/v1alpha3/role"
	"kube-aggregation/pkg/models/resources/v1alpha3/rolebinding"
	"kube-aggregation/pkg/models/resources/v1alpha3/secret"
	"kube-aggregation/pkg/models/resources/v1alpha3/service"
	"kube-aggregation/pkg/models/resources/v1alpha3/serviceaccount"
	"kube-aggregation/pkg/models/resources/v1alpha3/statefulset"
	"kube-aggregation/pkg/models/resources/v1alpha3/user"
	iamv1alpha2 "kubesphere.io/api/iam/v1alpha2"
)

var ErrResourceNotSupported = errors.New("resource is not supported")

type ResourceGetter struct {
	clusterResourceGetters    map[schema.GroupVersionResource]v1alpha3.Interface
	namespacedResourceGetters map[schema.GroupVersionResource]v1alpha3.Interface
}

func NewResourceGetter(factory informers.InformerFactory, cache cache.Cache) *ResourceGetter {
	namespacedResourceGetters := make(map[schema.GroupVersionResource]v1alpha3.Interface)
	clusterResourceGetters := make(map[schema.GroupVersionResource]v1alpha3.Interface)

	namespacedResourceGetters[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}] = deployment.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "daemonsets"}] = daemonset.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}] = statefulset.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}] = service.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}] = configmap.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "secrets"}] = secret.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}] = pod.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "serviceaccounts"}] = serviceaccount.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"}] = ingress.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"}] = networkpolicy.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"}] = job.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[schema.GroupVersionResource{Group: "app.k8s.io", Version: "v1beta1", Resource: "applications"}] = application.New(cache)
	clusterResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "persistentvolumes"}] = persistentvolume.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[rbacv1.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralRoleBinding)] = rolebinding.New(factory.KubernetesSharedInformerFactory())
	namespacedResourceGetters[rbacv1.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralRole)] = role.New(factory.KubernetesSharedInformerFactory())
	clusterResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "nodes"}] = node.New(factory.KubernetesSharedInformerFactory())
	clusterResourceGetters[schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}] = namespace.New(factory.KubernetesSharedInformerFactory())

	// kubesphere resources
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralGlobalRole)] = globalrole.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralUser)] = user.New(factory.KubeSphereSharedInformerFactory(), factory.KubernetesSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralGlobalRoleBinding)] = globalrolebinding.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralLoginRecord)] = loginrecord.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcePluralGroup)] = group.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[iamv1alpha2.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcePluralGroupBinding)] = groupbinding.New(factory.KubeSphereSharedInformerFactory())
	clusterResourceGetters[rbacv1.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralClusterRole)] = clusterrole.New(factory.KubernetesSharedInformerFactory())
	clusterResourceGetters[rbacv1.SchemeGroupVersion.WithResource(iamv1alpha2.ResourcesPluralClusterRoleBinding)] = clusterrolebinding.New(factory.KubernetesSharedInformerFactory())
	clusterResourceGetters[monitoringdashboardv1alpha2.GroupVersion.WithResource("clusterdashboards")] = clusterdashboard.New(cache)

	// federated resources
	namespacedResourceGetters[monitoringdashboardv1alpha2.GroupVersion.WithResource("dashboards")] = dashboard.New(cache)

	return &ResourceGetter{
		namespacedResourceGetters: namespacedResourceGetters,
		clusterResourceGetters:    clusterResourceGetters,
	}
}

// TryResource will retrieve a getter with resource name, it doesn't guarantee find resource with correct group version
// need to refactor this use schema.GroupVersionResource
func (r *ResourceGetter) TryResource(clusterScope bool, resource string) v1alpha3.Interface {
	if clusterScope {
		for k, v := range r.clusterResourceGetters {
			if k.Resource == resource {
				return v
			}
		}
	}
	for k, v := range r.namespacedResourceGetters {
		if k.Resource == resource {
			return v
		}
	}
	return nil
}

func (r *ResourceGetter) Get(resource, namespace, name string) (runtime.Object, error) {
	clusterScope := namespace == ""
	getter := r.TryResource(clusterScope, resource)
	if getter == nil {
		return nil, ErrResourceNotSupported
	}
	return getter.Get(namespace, name)
}

func (r *ResourceGetter) List(resource, namespace string, query *query.Query) (*api.ListResult, error) {
	clusterScope := namespace == ""
	getter := r.TryResource(clusterScope, resource)
	if getter == nil {
		return nil, ErrResourceNotSupported
	}
	return getter.List(namespace, query)
}
