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

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	typesv1beta1 "kubesphere.io/api/types/v1beta1"
	versioned "kube-aggregation/pkg/client/clientset/versioned"
	internalinterfaces "kube-aggregation/pkg/client/informers/externalversions/internalinterfaces"
	v1beta1 "kube-aggregation/pkg/client/listers/types/v1beta1"
)

// FederatedStatefulSetInformer provides access to a shared informer and lister for
// FederatedStatefulSets.
type FederatedStatefulSetInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.FederatedStatefulSetLister
}

type federatedStatefulSetInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewFederatedStatefulSetInformer constructs a new informer for FederatedStatefulSet type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFederatedStatefulSetInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredFederatedStatefulSetInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredFederatedStatefulSetInformer constructs a new informer for FederatedStatefulSet type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredFederatedStatefulSetInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TypesV1beta1().FederatedStatefulSets(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TypesV1beta1().FederatedStatefulSets(namespace).Watch(context.TODO(), options)
			},
		},
		&typesv1beta1.FederatedStatefulSet{},
		resyncPeriod,
		indexers,
	)
}

func (f *federatedStatefulSetInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredFederatedStatefulSetInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *federatedStatefulSetInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&typesv1beta1.FederatedStatefulSet{}, f.defaultInformer)
}

func (f *federatedStatefulSetInformer) Lister() v1beta1.FederatedStatefulSetLister {
	return v1beta1.NewFederatedStatefulSetLister(f.Informer().GetIndexer())
}
