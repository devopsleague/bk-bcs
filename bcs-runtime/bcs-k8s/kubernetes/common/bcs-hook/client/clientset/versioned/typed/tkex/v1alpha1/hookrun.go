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

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/kubernetes/common/bcs-hook/apis/tkex/v1alpha1"
	scheme "github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/kubernetes/common/bcs-hook/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// HookRunsGetter has a method to return a HookRunInterface.
// A group's client should implement this interface.
type HookRunsGetter interface {
	HookRuns(namespace string) HookRunInterface
}

// HookRunInterface has methods to work with HookRun resources.
type HookRunInterface interface {
	Create(context.Context, *v1alpha1.HookRun) (*v1alpha1.HookRun, error)
	Update(context.Context, *v1alpha1.HookRun) (*v1alpha1.HookRun, error)
	UpdateStatus(context.Context, *v1alpha1.HookRun) (*v1alpha1.HookRun, error)
	Delete(ctx context.Context, name string, options *v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(ctx context.Context, name string, options v1.GetOptions) (*v1alpha1.HookRun, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.HookRunList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.HookRun, err error)
	HookRunExpansion
}

// hookRuns implements HookRunInterface
type hookRuns struct {
	client rest.Interface
	ns     string
}

// newHookRuns returns a HookRuns
func newHookRuns(c *TkexV1alpha1Client, namespace string) *hookRuns {
	return &hookRuns{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the hookRun, and returns the corresponding hookRun object, and an error if there is any.
func (c *hookRuns) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.HookRun, err error) {
	result = &v1alpha1.HookRun{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("hookruns").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of HookRuns that match those selectors.
func (c *hookRuns) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.HookRunList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.HookRunList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("hookruns").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested hookRuns.
func (c *hookRuns) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("hookruns").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a hookRun and creates it.  Returns the server's representation of the hookRun, and an error, if there is any.
func (c *hookRuns) Create(ctx context.Context, hookRun *v1alpha1.HookRun) (result *v1alpha1.HookRun, err error) {
	result = &v1alpha1.HookRun{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("hookruns").
		Body(hookRun).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a hookRun and updates it. Returns the server's representation of the hookRun, and an error, if there is any.
func (c *hookRuns) Update(ctx context.Context, hookRun *v1alpha1.HookRun) (result *v1alpha1.HookRun, err error) {
	result = &v1alpha1.HookRun{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("hookruns").
		Name(hookRun.Name).
		Body(hookRun).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *hookRuns) UpdateStatus(ctx context.Context, hookRun *v1alpha1.HookRun) (result *v1alpha1.HookRun, err error) {
	result = &v1alpha1.HookRun{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("hookruns").
		Name(hookRun.Name).
		SubResource("status").
		Body(hookRun).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the hookRun and deletes it. Returns an error if one occurs.
func (c *hookRuns) Delete(ctx context.Context, name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("hookruns").
		Name(name).
		Body(options).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *hookRuns) DeleteCollection(ctx context.Context, options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("hookruns").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched hookRun.
func (c *hookRuns) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.HookRun, err error) {
	result = &v1alpha1.HookRun{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("hookruns").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
