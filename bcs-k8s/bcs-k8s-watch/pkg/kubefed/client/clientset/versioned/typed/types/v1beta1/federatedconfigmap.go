/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "bk-bcs/bcs-k8s/bcs-k8s-watch/pkg/kubefed/apis/types/v1beta1"
	scheme "bk-bcs/bcs-k8s/bcs-k8s-watch/pkg/kubefed/client/clientset/versioned/scheme"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// FederatedConfigMapsGetter has a method to return a FederatedConfigMapInterface.
// A group's client should implement this interface.
type FederatedConfigMapsGetter interface {
	FederatedConfigMaps(namespace string) FederatedConfigMapInterface
}

// FederatedConfigMapInterface has methods to work with FederatedConfigMap resources.
type FederatedConfigMapInterface interface {
	Create(*v1beta1.FederatedConfigMap) (*v1beta1.FederatedConfigMap, error)
	Update(*v1beta1.FederatedConfigMap) (*v1beta1.FederatedConfigMap, error)
	UpdateStatus(*v1beta1.FederatedConfigMap) (*v1beta1.FederatedConfigMap, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1beta1.FederatedConfigMap, error)
	List(opts v1.ListOptions) (*v1beta1.FederatedConfigMapList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.FederatedConfigMap, err error)
	FederatedConfigMapExpansion
}

// federatedConfigMaps implements FederatedConfigMapInterface
type federatedConfigMaps struct {
	client rest.Interface
	ns     string
}

// newFederatedConfigMaps returns a FederatedConfigMaps
func newFederatedConfigMaps(c *TypesV1beta1Client, namespace string) *federatedConfigMaps {
	return &federatedConfigMaps{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the federatedConfigMap, and returns the corresponding federatedConfigMap object, and an error if there is any.
func (c *federatedConfigMaps) Get(name string, options v1.GetOptions) (result *v1beta1.FederatedConfigMap, err error) {
	result = &v1beta1.FederatedConfigMap{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("federatedconfigmaps").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of FederatedConfigMaps that match those selectors.
func (c *federatedConfigMaps) List(opts v1.ListOptions) (result *v1beta1.FederatedConfigMapList, err error) {
	result = &v1beta1.FederatedConfigMapList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("federatedconfigmaps").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested federatedConfigMaps.
func (c *federatedConfigMaps) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("federatedconfigmaps").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a federatedConfigMap and creates it.  Returns the server's representation of the federatedConfigMap, and an error, if there is any.
func (c *federatedConfigMaps) Create(federatedConfigMap *v1beta1.FederatedConfigMap) (result *v1beta1.FederatedConfigMap, err error) {
	result = &v1beta1.FederatedConfigMap{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("federatedconfigmaps").
		Body(federatedConfigMap).
		Do().
		Into(result)
	return
}

// Update takes the representation of a federatedConfigMap and updates it. Returns the server's representation of the federatedConfigMap, and an error, if there is any.
func (c *federatedConfigMaps) Update(federatedConfigMap *v1beta1.FederatedConfigMap) (result *v1beta1.FederatedConfigMap, err error) {
	result = &v1beta1.FederatedConfigMap{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("federatedconfigmaps").
		Name(federatedConfigMap.Name).
		Body(federatedConfigMap).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *federatedConfigMaps) UpdateStatus(federatedConfigMap *v1beta1.FederatedConfigMap) (result *v1beta1.FederatedConfigMap, err error) {
	result = &v1beta1.FederatedConfigMap{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("federatedconfigmaps").
		Name(federatedConfigMap.Name).
		SubResource("status").
		Body(federatedConfigMap).
		Do().
		Into(result)
	return
}

// Delete takes name of the federatedConfigMap and deletes it. Returns an error if one occurs.
func (c *federatedConfigMaps) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("federatedconfigmaps").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *federatedConfigMaps) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("federatedconfigmaps").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched federatedConfigMap.
func (c *federatedConfigMaps) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.FederatedConfigMap, err error) {
	result = &v1beta1.FederatedConfigMap{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("federatedconfigmaps").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
