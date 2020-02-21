/*
Copyright 2020 The Kubernetes Authors.

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

package v2

import (
	v2 "bk-bcs/bcs-services/bcs-webhook-server/pkg/apis/bk-bcs/v2"
	scheme "bk-bcs/bcs-services/bcs-webhook-server/pkg/client/clientset/versioned/scheme"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// BcsLogConfigsGetter has a method to return a BcsLogConfigInterface.
// A group's client should implement this interface.
type BcsLogConfigsGetter interface {
	BcsLogConfigs(namespace string) BcsLogConfigInterface
}

// BcsLogConfigInterface has methods to work with BcsLogConfig resources.
type BcsLogConfigInterface interface {
	Create(*v2.BcsLogConfig) (*v2.BcsLogConfig, error)
	Update(*v2.BcsLogConfig) (*v2.BcsLogConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v2.BcsLogConfig, error)
	List(opts v1.ListOptions) (*v2.BcsLogConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.BcsLogConfig, err error)
	BcsLogConfigExpansion
}

// bcsLogConfigs implements BcsLogConfigInterface
type bcsLogConfigs struct {
	client rest.Interface
	ns     string
}

// newBcsLogConfigs returns a BcsLogConfigs
func newBcsLogConfigs(c *BkbcsV2Client, namespace string) *bcsLogConfigs {
	return &bcsLogConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the bcsLogConfig, and returns the corresponding bcsLogConfig object, and an error if there is any.
func (c *bcsLogConfigs) Get(name string, options v1.GetOptions) (result *v2.BcsLogConfig, err error) {
	result = &v2.BcsLogConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bcslogconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of BcsLogConfigs that match those selectors.
func (c *bcsLogConfigs) List(opts v1.ListOptions) (result *v2.BcsLogConfigList, err error) {
	result = &v2.BcsLogConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("bcslogconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested bcsLogConfigs.
func (c *bcsLogConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("bcslogconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a bcsLogConfig and creates it.  Returns the server's representation of the bcsLogConfig, and an error, if there is any.
func (c *bcsLogConfigs) Create(bcsLogConfig *v2.BcsLogConfig) (result *v2.BcsLogConfig, err error) {
	result = &v2.BcsLogConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("bcslogconfigs").
		Body(bcsLogConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a bcsLogConfig and updates it. Returns the server's representation of the bcsLogConfig, and an error, if there is any.
func (c *bcsLogConfigs) Update(bcsLogConfig *v2.BcsLogConfig) (result *v2.BcsLogConfig, err error) {
	result = &v2.BcsLogConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("bcslogconfigs").
		Name(bcsLogConfig.Name).
		Body(bcsLogConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the bcsLogConfig and deletes it. Returns an error if one occurs.
func (c *bcsLogConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bcslogconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *bcsLogConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("bcslogconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched bcsLogConfig.
func (c *bcsLogConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.BcsLogConfig, err error) {
	result = &v2.BcsLogConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("bcslogconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
