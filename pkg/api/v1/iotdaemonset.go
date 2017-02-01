package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const (
	IotDaemonSetKind = "IotDaemonSet"
	IotDaemonSetType = "iotdaemonsets"
)

type IotDaemonSet struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta       `json:"metadata,omitempty"`
	Spec            v1beta1.DaemonSetSpec   `json:"spec,omitempty"`
	Status          v1beta1.DaemonSetStatus `json:"status,omitempty"`
}

type IotDaemonSetList struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IotDaemonSet  `json:"items"`
}

func (iotDaemonSet *IotDaemonSet) GetObjectKind() schema.ObjectKind {
	return &iotDaemonSet.TypeMeta
}

func (iotDaemonSet *IotDaemonSet) GetObjectMeta() *metav1.ObjectMeta {
	return &iotDaemonSet.Metadata
}

func (iotDaemonSetList *IotDaemonSetList) GetObjectKind() schema.ObjectKind {
	return &iotDaemonSetList.TypeMeta
}

func (iotDaemonSetList *IotDaemonSetList) GetListMeta() metav1.List {
	return &iotDaemonSetList.Metadata
}
