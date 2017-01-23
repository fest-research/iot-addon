package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/pkg/api/v1"
)

const IotDeviceType = "iotdevices"

type IotDevice struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        v1.ObjectMeta `json:"metadata"`
	Spec            v1.NodeSpec   `json:"spec"`
	Status          v1.NodeStatus `json:"status"`
}

type IotDeviceList struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ListMeta `json:"metadata"`
	Items           []IotDevice     `json:"items"`
}

func (iotDevice *IotDevice) GetObjectKind() schema.ObjectKind {
	return &iotDevice.TypeMeta
}

func (iotDevice *IotDevice) GetObjectMeta() *v1.ObjectMeta {
	return &iotDevice.Metadata
}

func (iotDeviceList *IotDeviceList) GetObjectKind() schema.ObjectKind {
	return &iotDeviceList.TypeMeta
}

func (iotDeviceList *IotDeviceList) GetListMeta() metav1.List {
	return &iotDeviceList.Metadata
}
