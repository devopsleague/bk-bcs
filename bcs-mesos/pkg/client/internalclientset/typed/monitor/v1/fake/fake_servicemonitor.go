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

package fake

import (
	monitorv1 "github.com/Tencent/bk-bcs/bcs-mesos/pkg/apis/monitor/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeServiceMonitors implements ServiceMonitorInterface
type FakeServiceMonitors struct {
	Fake *FakeMonitorV1
	ns   string
}

var servicemonitorsResource = schema.GroupVersionResource{Group: "monitor.tencent.com", Version: "v1", Resource: "servicemonitors"}

var servicemonitorsKind = schema.GroupVersionKind{Group: "monitor.tencent.com", Version: "v1", Kind: "ServiceMonitor"}

// Get takes name of the serviceMonitor, and returns the corresponding serviceMonitor object, and an error if there is any.
func (c *FakeServiceMonitors) Get(name string, options v1.GetOptions) (result *monitorv1.ServiceMonitor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(servicemonitorsResource, c.ns, name), &monitorv1.ServiceMonitor{})

	if obj == nil {
		return nil, err
	}
	return obj.(*monitorv1.ServiceMonitor), err
}

// List takes label and field selectors, and returns the list of ServiceMonitors that match those selectors.
func (c *FakeServiceMonitors) List(opts v1.ListOptions) (result *monitorv1.ServiceMonitorList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(servicemonitorsResource, servicemonitorsKind, c.ns, opts), &monitorv1.ServiceMonitorList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &monitorv1.ServiceMonitorList{ListMeta: obj.(*monitorv1.ServiceMonitorList).ListMeta}
	for _, item := range obj.(*monitorv1.ServiceMonitorList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested serviceMonitors.
func (c *FakeServiceMonitors) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(servicemonitorsResource, c.ns, opts))

}

// Create takes the representation of a serviceMonitor and creates it.  Returns the server's representation of the serviceMonitor, and an error, if there is any.
func (c *FakeServiceMonitors) Create(serviceMonitor *monitorv1.ServiceMonitor) (result *monitorv1.ServiceMonitor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(servicemonitorsResource, c.ns, serviceMonitor), &monitorv1.ServiceMonitor{})

	if obj == nil {
		return nil, err
	}
	return obj.(*monitorv1.ServiceMonitor), err
}

// Update takes the representation of a serviceMonitor and updates it. Returns the server's representation of the serviceMonitor, and an error, if there is any.
func (c *FakeServiceMonitors) Update(serviceMonitor *monitorv1.ServiceMonitor) (result *monitorv1.ServiceMonitor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(servicemonitorsResource, c.ns, serviceMonitor), &monitorv1.ServiceMonitor{})

	if obj == nil {
		return nil, err
	}
	return obj.(*monitorv1.ServiceMonitor), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeServiceMonitors) UpdateStatus(serviceMonitor *monitorv1.ServiceMonitor) (*monitorv1.ServiceMonitor, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(servicemonitorsResource, "status", c.ns, serviceMonitor), &monitorv1.ServiceMonitor{})

	if obj == nil {
		return nil, err
	}
	return obj.(*monitorv1.ServiceMonitor), err
}

// Delete takes name of the serviceMonitor and deletes it. Returns an error if one occurs.
func (c *FakeServiceMonitors) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(servicemonitorsResource, c.ns, name), &monitorv1.ServiceMonitor{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServiceMonitors) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(servicemonitorsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &monitorv1.ServiceMonitorList{})
	return err
}

// Patch applies the patch and returns the patched serviceMonitor.
func (c *FakeServiceMonitors) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *monitorv1.ServiceMonitor, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(servicemonitorsResource, c.ns, name, data, subresources...), &monitorv1.ServiceMonitor{})

	if obj == nil {
		return nil, err
	}
	return obj.(*monitorv1.ServiceMonitor), err
}
