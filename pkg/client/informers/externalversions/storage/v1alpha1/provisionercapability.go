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

package v1alpha1

import (
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	storagev1alpha1 "kubesphere.io/api/storage/v1alpha1"
	versioned "kube-aggregation/pkg/client/clientset/versioned"
	internalinterfaces "kube-aggregation/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "kube-aggregation/pkg/client/listers/storage/v1alpha1"
)

// ProvisionerCapabilityInformer provides access to a shared informer and lister for
// ProvisionerCapabilities.
type ProvisionerCapabilityInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ProvisionerCapabilityLister
}

type provisionerCapabilityInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewProvisionerCapabilityInformer constructs a new informer for ProvisionerCapability type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewProvisionerCapabilityInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredProvisionerCapabilityInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredProvisionerCapabilityInformer constructs a new informer for ProvisionerCapability type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredProvisionerCapabilityInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.StorageV1alpha1().ProvisionerCapabilities().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.StorageV1alpha1().ProvisionerCapabilities().Watch(context.TODO(), options)
			},
		},
		&storagev1alpha1.ProvisionerCapability{},
		resyncPeriod,
		indexers,
	)
}

func (f *provisionerCapabilityInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredProvisionerCapabilityInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *provisionerCapabilityInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&storagev1alpha1.ProvisionerCapability{}, f.defaultInformer)
}

func (f *provisionerCapabilityInformer) Lister() v1alpha1.ProvisionerCapabilityLister {
	return v1alpha1.NewProvisionerCapabilityLister(f.Informer().GetIndexer())
}
