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
	applicationv1alpha1 "kubesphere.io/api/application/v1alpha1"
	versioned "kube-aggregation/pkg/client/clientset/versioned"
	internalinterfaces "kube-aggregation/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "kube-aggregation/pkg/client/listers/application/v1alpha1"
)

// HelmApplicationVersionInformer provides access to a shared informer and lister for
// HelmApplicationVersions.
type HelmApplicationVersionInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.HelmApplicationVersionLister
}

type helmApplicationVersionInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewHelmApplicationVersionInformer constructs a new informer for HelmApplicationVersion type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewHelmApplicationVersionInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredHelmApplicationVersionInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredHelmApplicationVersionInformer constructs a new informer for HelmApplicationVersion type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredHelmApplicationVersionInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApplicationV1alpha1().HelmApplicationVersions().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApplicationV1alpha1().HelmApplicationVersions().Watch(context.TODO(), options)
			},
		},
		&applicationv1alpha1.HelmApplicationVersion{},
		resyncPeriod,
		indexers,
	)
}

func (f *helmApplicationVersionInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredHelmApplicationVersionInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *helmApplicationVersionInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&applicationv1alpha1.HelmApplicationVersion{}, f.defaultInformer)
}

func (f *helmApplicationVersionInformer) Lister() v1alpha1.HelmApplicationVersionLister {
	return v1alpha1.NewHelmApplicationVersionLister(f.Informer().GetIndexer())
}
