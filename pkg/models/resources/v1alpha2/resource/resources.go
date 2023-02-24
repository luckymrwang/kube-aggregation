/*
Copyright 2019 The KubeAggregation Authors.

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

	"k8s.io/klog"

	"github.com/inspur/pkg/apiserver/params"
	"github.com/inspur/pkg/informers"
	"github.com/inspur/pkg/models"
	"github.com/inspur/pkg/models/resources/v1alpha2"
	"github.com/inspur/pkg/models/resources/v1alpha2/clusterrole"
	"github.com/inspur/pkg/models/resources/v1alpha2/configmap"
	"github.com/inspur/pkg/models/resources/v1alpha2/cronjob"
	"github.com/inspur/pkg/models/resources/v1alpha2/daemonset"
	"github.com/inspur/pkg/models/resources/v1alpha2/deployment"
	"github.com/inspur/pkg/models/resources/v1alpha2/hpa"
	"github.com/inspur/pkg/models/resources/v1alpha2/ingress"
	"github.com/inspur/pkg/models/resources/v1alpha2/job"
	"github.com/inspur/pkg/models/resources/v1alpha2/namespace"
	"github.com/inspur/pkg/models/resources/v1alpha2/node"
	"github.com/inspur/pkg/models/resources/v1alpha2/pod"
	"github.com/inspur/pkg/models/resources/v1alpha2/role"
	"github.com/inspur/pkg/models/resources/v1alpha2/secret"
	"github.com/inspur/pkg/models/resources/v1alpha2/service"
	"github.com/inspur/pkg/models/resources/v1alpha2/statefulset"
	"github.com/inspur/pkg/utils/sliceutil"
)

var ErrResourceNotSupported = errors.New("resource is not supported")

type ResourceGetter struct {
	resourcesGetters map[string]v1alpha2.Interface
}

func (r ResourceGetter) Add(resource string, getter v1alpha2.Interface) {
	if r.resourcesGetters == nil {
		r.resourcesGetters = make(map[string]v1alpha2.Interface)
	}
	r.resourcesGetters[resource] = getter
}

func NewResourceGetter(factory informers.InformerFactory) *ResourceGetter {
	resourceGetters := make(map[string]v1alpha2.Interface)

	resourceGetters[v1alpha2.ConfigMaps] = configmap.NewConfigmapSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.CronJobs] = cronjob.NewCronJobSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.DaemonSets] = daemonset.NewDaemonSetSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Deployments] = deployment.NewDeploymentSetSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Ingresses] = ingress.NewIngressSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Jobs] = job.NewJobSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Secrets] = secret.NewSecretSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Services] = service.NewServiceSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.StatefulSets] = statefulset.NewStatefulSetSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Pods] = pod.NewPodSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Roles] = role.NewRoleSearcher(factory.KubernetesSharedInformerFactory())

	resourceGetters[v1alpha2.Nodes] = node.NewNodeSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.Namespaces] = namespace.NewNamespaceSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.ClusterRoles] = clusterrole.NewClusterRoleSearcher(factory.KubernetesSharedInformerFactory())
	resourceGetters[v1alpha2.HorizontalPodAutoscalers] = hpa.NewHpaSearcher(factory.KubernetesSharedInformerFactory())

	return &ResourceGetter{resourcesGetters: resourceGetters}

}

var (
	clusterResources = []string{v1alpha2.Nodes, v1alpha2.Workspaces, v1alpha2.Namespaces, v1alpha2.ClusterRoles, v1alpha2.StorageClasses, v1alpha2.S2iBuilderTemplates}
)

func (r *ResourceGetter) GetResource(namespace, resource, name string) (interface{}, error) {
	// none namespace resource
	if namespace != "" && sliceutil.HasString(clusterResources, resource) {
		return nil, ErrResourceNotSupported
	}
	if searcher, ok := r.resourcesGetters[resource]; ok {
		resource, err := searcher.Get(namespace, name)
		if err != nil {
			klog.Error(err)
			return nil, err
		}
		return resource, nil
	}
	return nil, ErrResourceNotSupported
}

func (r *ResourceGetter) ListResources(namespace, resource string, conditions *params.Conditions, orderBy string, reverse bool, limit, offset int) (*models.PageableResponse, error) {
	items := make([]interface{}, 0)
	var err error
	var result []interface{}

	// none namespace resource
	if namespace != "" && sliceutil.HasString(clusterResources, resource) {
		return nil, ErrResourceNotSupported
	}

	if searcher, ok := r.resourcesGetters[resource]; ok {
		result, err = searcher.Search(namespace, conditions, orderBy, reverse)
	} else {
		return nil, ErrResourceNotSupported
	}

	if err != nil {
		klog.Error(err)
		return nil, err
	}

	if limit == -1 || limit+offset > len(result) {
		limit = len(result) - offset
	}

	items = result[offset : offset+limit]

	return &models.PageableResponse{TotalCount: len(result), Items: items}, nil
}
