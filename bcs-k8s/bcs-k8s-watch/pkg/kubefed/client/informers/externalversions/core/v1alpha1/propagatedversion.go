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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	core_v1alpha1 "bk-bcs/bcs-k8s/bcs-k8s-watch/pkg/kubefed/apis/core/v1alpha1"
	versioned "bk-bcs/bcs-k8s/bcs-k8s-watch/pkg/kubefed/client/clientset/versioned"
	internalinterfaces "bk-bcs/bcs-k8s/bcs-k8s-watch/pkg/kubefed/client/informers/externalversions/internalinterfaces"
	v1alpha1 "bk-bcs/bcs-k8s/bcs-k8s-watch/pkg/kubefed/client/listers/core/v1alpha1"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// PropagatedVersionInformer provides access to a shared informer and lister for
// PropagatedVersions.
type PropagatedVersionInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.PropagatedVersionLister
}

type propagatedVersionInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewPropagatedVersionInformer constructs a new informer for PropagatedVersion type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewPropagatedVersionInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredPropagatedVersionInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredPropagatedVersionInformer constructs a new informer for PropagatedVersion type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredPropagatedVersionInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1alpha1().PropagatedVersions(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CoreV1alpha1().PropagatedVersions(namespace).Watch(options)
			},
		},
		&core_v1alpha1.PropagatedVersion{},
		resyncPeriod,
		indexers,
	)
}

func (f *propagatedVersionInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredPropagatedVersionInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *propagatedVersionInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&core_v1alpha1.PropagatedVersion{}, f.defaultInformer)
}

func (f *propagatedVersionInformer) Lister() v1alpha1.PropagatedVersionLister {
	return v1alpha1.NewPropagatedVersionLister(f.Informer().GetIndexer())
}
