/*
Copyright The Kubernetes Authors.

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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	crdexamplecomv1 "github.com/murali-bashyam/crdexample/pkg/apis/crd.example.com/v1"
	versioned "github.com/murali-bashyam/crdexample/pkg/client/clientset/versioned"
	internalinterfaces "github.com/murali-bashyam/crdexample/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/murali-bashyam/crdexample/pkg/client/listers/crd.example.com/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// StoragePoolInformer provides access to a shared informer and lister for
// StoragePools.
type StoragePoolInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.StoragePoolLister
}

type storagePoolInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewStoragePoolInformer constructs a new informer for StoragePool type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewStoragePoolInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredStoragePoolInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredStoragePoolInformer constructs a new informer for StoragePool type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredStoragePoolInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CrdV1().StoragePools(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CrdV1().StoragePools(namespace).Watch(options)
			},
		},
		&crdexamplecomv1.StoragePool{},
		resyncPeriod,
		indexers,
	)
}

func (f *storagePoolInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredStoragePoolInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *storagePoolInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&crdexamplecomv1.StoragePool{}, f.defaultInformer)
}

func (f *storagePoolInformer) Lister() v1.StoragePoolLister {
	return v1.NewStoragePoolLister(f.Informer().GetIndexer())
}
