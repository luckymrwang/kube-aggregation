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

package informers

import (
	"reflect"
	"time"

	prominformers "github.com/prometheus-operator/prometheus-operator/pkg/client/informers/externalversions"
	k8sinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"

	"kube-aggregation/pkg/client/clientset/versioned"
	ksinformers "kube-aggregation/pkg/client/informers/externalversions"
)

// default re-sync period for all informer factories
const defaultResync = 600 * time.Second

// InformerFactory is a group all shared informer factories which kubesphere needed
// callers should check if the return value is nil
type InformerFactory interface {
	KubernetesSharedInformerFactory() k8sinformers.SharedInformerFactory
	KubeSphereSharedInformerFactory() ksinformers.SharedInformerFactory

	// Start shared informer factory one by one if they are not nil
	Start(stopCh <-chan struct{})
}

type GenericInformerFactory interface {
	Start(stopCh <-chan struct{})
	WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool
}

type informerFactories struct {
	informerFactory           k8sinformers.SharedInformerFactory
	ksInformerFactory         ksinformers.SharedInformerFactory
	prometheusInformerFactory prominformers.SharedInformerFactory
}

func NewInformerFactories(client kubernetes.Interface, ksClient versioned.Interface) InformerFactory {
	factory := &informerFactories{}

	if client != nil {
		factory.informerFactory = k8sinformers.NewSharedInformerFactory(client, defaultResync)
	}

	if ksClient != nil {
		factory.ksInformerFactory = ksinformers.NewSharedInformerFactory(ksClient, defaultResync)
	}

	return factory
}

func (f *informerFactories) KubernetesSharedInformerFactory() k8sinformers.SharedInformerFactory {
	return f.informerFactory
}

func (f *informerFactories) KubeSphereSharedInformerFactory() ksinformers.SharedInformerFactory {
	return f.ksInformerFactory
}

func (f *informerFactories) PrometheusSharedInformerFactory() prominformers.SharedInformerFactory {
	return f.prometheusInformerFactory
}

func (f *informerFactories) Start(stopCh <-chan struct{}) {
	if f.informerFactory != nil {
		f.informerFactory.Start(stopCh)
	}

	if f.ksInformerFactory != nil {
		f.ksInformerFactory.Start(stopCh)
	}

	if f.prometheusInformerFactory != nil {
		f.prometheusInformerFactory.Start(stopCh)
	}
}
