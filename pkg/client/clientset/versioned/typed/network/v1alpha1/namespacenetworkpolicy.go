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

package v1alpha1

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "kubesphere.io/api/network/v1alpha1"
	scheme "kube-aggregation/pkg/client/clientset/versioned/scheme"
)

// NamespaceNetworkPoliciesGetter has a method to return a NamespaceNetworkPolicyInterface.
// A group's client should implement this interface.
type NamespaceNetworkPoliciesGetter interface {
	NamespaceNetworkPolicies(namespace string) NamespaceNetworkPolicyInterface
}

// NamespaceNetworkPolicyInterface has methods to work with NamespaceNetworkPolicy resources.
type NamespaceNetworkPolicyInterface interface {
	Create(ctx context.Context, namespaceNetworkPolicy *v1alpha1.NamespaceNetworkPolicy, opts v1.CreateOptions) (*v1alpha1.NamespaceNetworkPolicy, error)
	Update(ctx context.Context, namespaceNetworkPolicy *v1alpha1.NamespaceNetworkPolicy, opts v1.UpdateOptions) (*v1alpha1.NamespaceNetworkPolicy, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.NamespaceNetworkPolicy, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.NamespaceNetworkPolicyList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.NamespaceNetworkPolicy, err error)
	NamespaceNetworkPolicyExpansion
}

// namespaceNetworkPolicies implements NamespaceNetworkPolicyInterface
type namespaceNetworkPolicies struct {
	client rest.Interface
	ns     string
}

// newNamespaceNetworkPolicies returns a NamespaceNetworkPolicies
func newNamespaceNetworkPolicies(c *NetworkV1alpha1Client, namespace string) *namespaceNetworkPolicies {
	return &namespaceNetworkPolicies{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the namespaceNetworkPolicy, and returns the corresponding namespaceNetworkPolicy object, and an error if there is any.
func (c *namespaceNetworkPolicies) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.NamespaceNetworkPolicy, err error) {
	result = &v1alpha1.NamespaceNetworkPolicy{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("namespacenetworkpolicies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of NamespaceNetworkPolicies that match those selectors.
func (c *namespaceNetworkPolicies) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.NamespaceNetworkPolicyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.NamespaceNetworkPolicyList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("namespacenetworkpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested namespaceNetworkPolicies.
func (c *namespaceNetworkPolicies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("namespacenetworkpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a namespaceNetworkPolicy and creates it.  Returns the server's representation of the namespaceNetworkPolicy, and an error, if there is any.
func (c *namespaceNetworkPolicies) Create(ctx context.Context, namespaceNetworkPolicy *v1alpha1.NamespaceNetworkPolicy, opts v1.CreateOptions) (result *v1alpha1.NamespaceNetworkPolicy, err error) {
	result = &v1alpha1.NamespaceNetworkPolicy{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("namespacenetworkpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(namespaceNetworkPolicy).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a namespaceNetworkPolicy and updates it. Returns the server's representation of the namespaceNetworkPolicy, and an error, if there is any.
func (c *namespaceNetworkPolicies) Update(ctx context.Context, namespaceNetworkPolicy *v1alpha1.NamespaceNetworkPolicy, opts v1.UpdateOptions) (result *v1alpha1.NamespaceNetworkPolicy, err error) {
	result = &v1alpha1.NamespaceNetworkPolicy{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("namespacenetworkpolicies").
		Name(namespaceNetworkPolicy.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(namespaceNetworkPolicy).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the namespaceNetworkPolicy and deletes it. Returns an error if one occurs.
func (c *namespaceNetworkPolicies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("namespacenetworkpolicies").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *namespaceNetworkPolicies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("namespacenetworkpolicies").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched namespaceNetworkPolicy.
func (c *namespaceNetworkPolicies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.NamespaceNetworkPolicy, err error) {
	result = &v1alpha1.NamespaceNetworkPolicy{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("namespacenetworkpolicies").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
