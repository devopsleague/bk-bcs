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

// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "bk-bcs/bcs-k8s/bcs-k8s-watch/pkg/kubefed/apis/types/v1beta1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// FederatedRoleBindingLister helps list FederatedRoleBindings.
type FederatedRoleBindingLister interface {
	// List lists all FederatedRoleBindings in the indexer.
	List(selector labels.Selector) (ret []*v1beta1.FederatedRoleBinding, err error)
	// FederatedRoleBindings returns an object that can list and get FederatedRoleBindings.
	FederatedRoleBindings(namespace string) FederatedRoleBindingNamespaceLister
	FederatedRoleBindingListerExpansion
}

// federatedRoleBindingLister implements the FederatedRoleBindingLister interface.
type federatedRoleBindingLister struct {
	indexer cache.Indexer
}

// NewFederatedRoleBindingLister returns a new FederatedRoleBindingLister.
func NewFederatedRoleBindingLister(indexer cache.Indexer) FederatedRoleBindingLister {
	return &federatedRoleBindingLister{indexer: indexer}
}

// List lists all FederatedRoleBindings in the indexer.
func (s *federatedRoleBindingLister) List(selector labels.Selector) (ret []*v1beta1.FederatedRoleBinding, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.FederatedRoleBinding))
	})
	return ret, err
}

// FederatedRoleBindings returns an object that can list and get FederatedRoleBindings.
func (s *federatedRoleBindingLister) FederatedRoleBindings(namespace string) FederatedRoleBindingNamespaceLister {
	return federatedRoleBindingNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// FederatedRoleBindingNamespaceLister helps list and get FederatedRoleBindings.
type FederatedRoleBindingNamespaceLister interface {
	// List lists all FederatedRoleBindings in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1beta1.FederatedRoleBinding, err error)
	// Get retrieves the FederatedRoleBinding from the indexer for a given namespace and name.
	Get(name string) (*v1beta1.FederatedRoleBinding, error)
	FederatedRoleBindingNamespaceListerExpansion
}

// federatedRoleBindingNamespaceLister implements the FederatedRoleBindingNamespaceLister
// interface.
type federatedRoleBindingNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all FederatedRoleBindings in the indexer for a given namespace.
func (s federatedRoleBindingNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.FederatedRoleBinding, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.FederatedRoleBinding))
	})
	return ret, err
}

// Get retrieves the FederatedRoleBinding from the indexer for a given namespace and name.
func (s federatedRoleBindingNamespaceLister) Get(name string) (*v1beta1.FederatedRoleBinding, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("federatedrolebinding"), name)
	}
	return obj.(*v1beta1.FederatedRoleBinding), nil
}
