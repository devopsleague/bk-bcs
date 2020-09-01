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
	v2 "github.com/Tencent/bk-bcs/bcs-mesos/pkg/apis/bkbcs/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBcsCommandInfos implements BcsCommandInfoInterface
type FakeBcsCommandInfos struct {
	Fake *FakeBkbcsV2
	ns   string
}

var bcscommandinfosResource = schema.GroupVersionResource{Group: "bkbcs.tencent.com", Version: "v2", Resource: "bcscommandinfos"}

var bcscommandinfosKind = schema.GroupVersionKind{Group: "bkbcs.tencent.com", Version: "v2", Kind: "BcsCommandInfo"}

// Get takes name of the bcsCommandInfo, and returns the corresponding bcsCommandInfo object, and an error if there is any.
func (c *FakeBcsCommandInfos) Get(name string, options v1.GetOptions) (result *v2.BcsCommandInfo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(bcscommandinfosResource, c.ns, name), &v2.BcsCommandInfo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.BcsCommandInfo), err
}

// List takes label and field selectors, and returns the list of BcsCommandInfos that match those selectors.
func (c *FakeBcsCommandInfos) List(opts v1.ListOptions) (result *v2.BcsCommandInfoList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(bcscommandinfosResource, bcscommandinfosKind, c.ns, opts), &v2.BcsCommandInfoList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v2.BcsCommandInfoList{ListMeta: obj.(*v2.BcsCommandInfoList).ListMeta}
	for _, item := range obj.(*v2.BcsCommandInfoList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested bcsCommandInfos.
func (c *FakeBcsCommandInfos) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(bcscommandinfosResource, c.ns, opts))

}

// Create takes the representation of a bcsCommandInfo and creates it.  Returns the server's representation of the bcsCommandInfo, and an error, if there is any.
func (c *FakeBcsCommandInfos) Create(bcsCommandInfo *v2.BcsCommandInfo) (result *v2.BcsCommandInfo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(bcscommandinfosResource, c.ns, bcsCommandInfo), &v2.BcsCommandInfo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.BcsCommandInfo), err
}

// Update takes the representation of a bcsCommandInfo and updates it. Returns the server's representation of the bcsCommandInfo, and an error, if there is any.
func (c *FakeBcsCommandInfos) Update(bcsCommandInfo *v2.BcsCommandInfo) (result *v2.BcsCommandInfo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(bcscommandinfosResource, c.ns, bcsCommandInfo), &v2.BcsCommandInfo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.BcsCommandInfo), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBcsCommandInfos) UpdateStatus(bcsCommandInfo *v2.BcsCommandInfo) (*v2.BcsCommandInfo, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(bcscommandinfosResource, "status", c.ns, bcsCommandInfo), &v2.BcsCommandInfo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.BcsCommandInfo), err
}

// Delete takes name of the bcsCommandInfo and deletes it. Returns an error if one occurs.
func (c *FakeBcsCommandInfos) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(bcscommandinfosResource, c.ns, name), &v2.BcsCommandInfo{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBcsCommandInfos) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(bcscommandinfosResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v2.BcsCommandInfoList{})
	return err
}

// Patch applies the patch and returns the patched bcsCommandInfo.
func (c *FakeBcsCommandInfos) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.BcsCommandInfo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(bcscommandinfosResource, c.ns, name, data, subresources...), &v2.BcsCommandInfo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.BcsCommandInfo), err
}
