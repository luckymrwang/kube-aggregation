/*
Copyright 2020 KubeSphere Authors

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

package v1alpha2

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"

	"kube-aggregation/pkg/api"
	"kube-aggregation/pkg/informers"
	"kube-aggregation/pkg/models/components"
	"kube-aggregation/pkg/models/kubeconfig"
	"kube-aggregation/pkg/models/kubectl"
	"kube-aggregation/pkg/models/resources/v1alpha2"
	"kube-aggregation/pkg/models/resources/v1alpha2/resource"
	"kube-aggregation/pkg/models/revisions"
	"kube-aggregation/pkg/server/errors"
	"kube-aggregation/pkg/server/params"
)

type resourceHandler struct {
	resourcesGetter    *resource.ResourceGetter
	componentsGetter   components.ComponentsGetter
	revisionGetter     revisions.RevisionGetter
	kubeconfigOperator kubeconfig.Interface
	kubectlOperator    kubectl.Interface
}

func newResourceHandler(k8sClient kubernetes.Interface, factory informers.InformerFactory, masterURL string) *resourceHandler {

	return &resourceHandler{
		resourcesGetter:    resource.NewResourceGetter(factory),
		componentsGetter:   components.NewComponentsGetter(factory.KubernetesSharedInformerFactory()),
		revisionGetter:     revisions.NewRevisionGetter(factory.KubernetesSharedInformerFactory()),
		kubeconfigOperator: kubeconfig.NewReadOnlyOperator(factory.KubernetesSharedInformerFactory().Core().V1().ConfigMaps().Lister(), masterURL),
		kubectlOperator: kubectl.NewOperator(nil, factory.KubernetesSharedInformerFactory().Apps().V1().Deployments(),
			factory.KubernetesSharedInformerFactory().Core().V1().Pods(), ""),
	}
}

func (r *resourceHandler) handleGetNamespacedResources(request *restful.Request, response *restful.Response) {
	r.handleListNamespaceResources(request, response)
}

func (r *resourceHandler) handleListNamespaceResources(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	resource := request.PathParameter("resources")
	orderBy := params.GetStringValueWithDefault(request, params.OrderByParam, v1alpha2.CreateTime)
	limit, offset := params.ParsePaging(request)
	reverse := params.GetBoolValueWithDefault(request, params.ReverseParam, false)
	conditions, err := params.ParseConditions(request)

	if err != nil {
		klog.Error(err)
		api.HandleBadRequest(response, request, err)
		return
	}

	result, err := r.resourcesGetter.ListResources(namespace, resource, conditions, orderBy, reverse, limit, offset)

	if err != nil {
		klog.Error(err)
		api.HandleInternalError(response, nil, err)
		return
	}

	response.WriteEntity(result)
}

func (r *resourceHandler) handleGetSystemHealthStatus(_ *restful.Request, response *restful.Response) {
	result, err := r.componentsGetter.GetSystemHealthStatus()

	if err != nil {
		api.HandleInternalError(response, nil, err)
		return
	}

	response.WriteAsJson(result)
}

func (r *resourceHandler) handleGetComponentStatus(request *restful.Request, response *restful.Response) {
	component := request.PathParameter("component")
	result, err := r.componentsGetter.GetComponentStatus(component)

	if err != nil {
		api.HandleInternalError(response, nil, err)
		return
	}

	response.WriteAsJson(result)
}

func (r *resourceHandler) handleGetComponents(_ *restful.Request, response *restful.Response) {
	result, err := r.componentsGetter.GetAllComponentsStatus()

	if err != nil {
		api.HandleInternalError(response, nil, err)
		return
	}

	response.WriteAsJson(result)
}

func (r *resourceHandler) handleGetDaemonSetRevision(request *restful.Request, response *restful.Response) {
	daemonset := request.PathParameter("daemonset")
	namespace := request.PathParameter("namespace")
	revision, err := strconv.Atoi(request.PathParameter("revision"))

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusBadRequest, errors.Wrap(err))
		return
	}

	result, err := r.revisionGetter.GetDaemonSetRevision(namespace, daemonset, revision)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, errors.Wrap(err))
		return
	}

	response.WriteAsJson(result)
}

func (r *resourceHandler) handleGetDeploymentRevision(request *restful.Request, response *restful.Response) {
	deploy := request.PathParameter("deployment")
	namespace := request.PathParameter("namespace")
	revision := request.PathParameter("revision")

	result, err := r.revisionGetter.GetDeploymentRevision(namespace, deploy, revision)

	if err != nil {
		response.WriteHeaderAndEntity(http.StatusInternalServerError, errors.Wrap(err))
		return
	}

	response.WriteAsJson(result)
}

func (r *resourceHandler) handleGetStatefulSetRevision(request *restful.Request, response *restful.Response) {
	statefulset := request.PathParameter("statefulset")
	namespace := request.PathParameter("namespace")
	revision, err := strconv.Atoi(request.PathParameter("revision"))
	if err != nil {
		response.WriteHeaderAndEntity(http.StatusBadRequest, errors.Wrap(err))
		return
	}

	result, err := r.revisionGetter.GetStatefulSetRevision(namespace, statefulset, revision)
	if err != nil {
		api.HandleInternalError(response, nil, err)
		return
	}
	response.WriteAsJson(result)
}

func (r *resourceHandler) handleGetNamespacedAbnormalWorkloads(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")

	result := api.Workloads{
		Namespace: namespace,
		Count:     make(map[string]int),
	}

	for _, workloadType := range []string{api.ResourceKindDeployment, api.ResourceKindStatefulSet, api.ResourceKindDaemonSet, api.ResourceKindJob, api.ResourceKindPersistentVolumeClaim} {
		var notReadyStatus string

		switch workloadType {
		case api.ResourceKindPersistentVolumeClaim:
			notReadyStatus = strings.Join([]string{v1alpha2.StatusPending, v1alpha2.StatusLost}, "|")
		case api.ResourceKindJob:
			notReadyStatus = v1alpha2.StatusFailed
		default:
			notReadyStatus = v1alpha2.StatusUpdating
		}

		res, err := r.resourcesGetter.ListResources(namespace, workloadType, &params.Conditions{Match: map[string]string{v1alpha2.Status: notReadyStatus}}, "", false, -1, 0)
		if err != nil {
			api.HandleInternalError(response, nil, err)
		}

		result.Count[workloadType] = len(res.Items)
	}

	response.WriteAsJson(result)

}

func (r *resourceHandler) GetKubectlPod(request *restful.Request, response *restful.Response) {
	user := request.PathParameter("user")

	kubectlPod, err := r.kubectlOperator.GetKubectlPod(user)

	if err != nil {
		klog.Errorln(err)
		response.WriteHeaderAndEntity(http.StatusInternalServerError, errors.Wrap(err))
		return
	}

	response.WriteEntity(kubectlPod)
}

func (r *resourceHandler) GetKubeconfig(request *restful.Request, response *restful.Response) {
	user := request.PathParameter("user")

	kubectlConfig, err := r.kubeconfigOperator.GetKubeConfig(user)

	if err != nil {
		klog.Error(err)
		if k8serr.IsNotFound(err) {
			// recreate
			response.WriteHeaderAndJson(http.StatusNotFound, errors.Wrap(err), restful.MIME_JSON)
		} else {
			response.WriteHeaderAndJson(http.StatusInternalServerError, errors.Wrap(err), restful.MIME_JSON)
		}
		return
	}

	response.Write([]byte(kubectlConfig))
}
