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

// FederatedNamespaceInformer provides access to a shared informer and lister for
// FederatedNamespaces.
type FederatedNamespaceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.FederatedNamespaceLister
}

type federatedNamespaceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewFederatedNamespaceInformer constructs a new informer for FederatedNamespace type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFederatedNamespaceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredFederatedNamespaceInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredFederatedNamespaceInformer constructs a new informer for FederatedNamespace type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredFederatedNamespaceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TypesV1beta1().FederatedNamespaces(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TypesV1beta1().FederatedNamespaces(namespace).Watch(context.TODO(), options)
			},
		},
		&typesv1beta1.FederatedNamespace{},
		resyncPeriod,
		indexers,
	)
}

func (f *federatedNamespaceInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredFederatedNamespaceInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *federatedNamespaceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&typesv1beta1.FederatedNamespace{}, f.defaultInformer)
}

func (f *federatedNamespaceInformer) Lister() v1beta1.FederatedNamespaceLister {
	return v1beta1.NewFederatedNamespaceLister(f.Informer().GetIndexer())
}
