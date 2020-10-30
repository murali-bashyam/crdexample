package v1

import (
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// StoragePool describes a storage pool.
 type StoragePool struct {
 metav1.TypeMeta `json:",inline"`
 metav1.ObjectMeta `json:"metadata,omitempty"`


 PoolSpec StoragePoolSpec `json:"spec"`
 }


// StoragePoolSpec is the spec for a storage pool resource
 type StoragePoolSpec struct {
 FailureDomain string `json:"failureDomain"`
 Quota int `json:"quota"`
 }


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object


// StoragePoolList is a list of storage pool resources
 type StoragePoolList struct {
 metav1.TypeMeta `json:",inline"`
 metav1.ListMeta `json:"metadata"`


 Items []StoragePool `json:"items"`
 }
