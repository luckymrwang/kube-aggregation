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

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1beta1 "kubesphere.io/api/types/v1beta1"
	scheme "kube-aggregation/pkg/client/clientset/versioned/scheme"
)

// FederatedServicesGetter has a method to return a FederatedServiceInterface.
// A group's client should implement this interface.
type FederatedServicesGetter interface {
	FederatedServices(namespace string) FederatedServiceInterface
}

// FederatedServiceInterface has methods to work with FederatedService resources.
type FederatedServiceInterface interface {
	Create(ctx context.Context, federatedService *v1beta1.FederatedService, opts v1.CreateOptions) (*v1beta1.FederatedService, error)
	Update(ctx context.Context, federatedService *v1beta1.FederatedService, opts v1.UpdateOptions) (*v1beta1.FederatedService, error)
	UpdateStatus(ctx context.Context, federatedService *v1beta1.FederatedService, opts v1.UpdateOptions) (*v1beta1.FederatedService, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.FederatedService, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.FederatedServiceList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.FederatedService, err error)
	FederatedServiceExpansion
}

// federatedServices implements FederatedServiceInterface
type federatedServices struct {
	client rest.Interface
	ns     string
}

// newFederatedServices returns a FederatedServices
func newFederatedServices(c *TypesV1beta1Client, namespace string) *federatedServices {
	return &federatedServices{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the federatedService, and returns the corresponding federatedService object, and an error if there is any.
func (c *federatedServices) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.FederatedService, err error) {
	result = &v1beta1.FederatedService{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("federatedservices").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FederatedServices that match those selectors.
func (c *federatedServices) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.FederatedServiceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.FederatedServiceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("federatedservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested federatedServices.
func (c *federatedServices) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("federatedservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a federatedService and creates it.  Returns the server's representation of the federatedService, and an error, if there is any.
func (c *federatedServices) Create(ctx context.Context, federatedService *v1beta1.FederatedService, opts v1.CreateOptions) (result *v1beta1.FederatedService, err error) {
	result = &v1beta1.FederatedService{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("federatedservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(federatedService).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a federatedService and updates it. Returns the server's representation of the federatedService, and an error, if there is any.
func (c *federatedServices) Update(ctx context.Context, federatedService *v1beta1.FederatedService, opts v1.UpdateOptions) (result *v1beta1.FederatedService, err error) {
	result = &v1beta1.FederatedService{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("federatedservices").
		Name(federatedService.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(federatedService).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *federatedServices) UpdateStatus(ctx context.Context, federatedService *v1beta1.FederatedService, opts v1.UpdateOptions) (result *v1beta1.FederatedService, err error) {
	result = &v1beta1.FederatedService{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("federatedservices").
		Name(federatedService.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(federatedService).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the federatedService and deletes it. Returns an error if one occurs.
func (c *federatedServices) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("federatedservices").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *federatedServices) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("federatedservices").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched federatedService.
func (c *federatedServices) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.FederatedService, err error) {
	result = &v1beta1.FederatedService{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("federatedservices").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
