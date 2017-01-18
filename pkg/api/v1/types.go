package v1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type IotDevice struct {
	v1.TypeMeta `json:",inline"`
	Metadata    v1.ListMeta `json:"metadata"`
}

type IotDeviceList struct {
	v1.TypeMeta `json:",inline"`
	Metadata    v1.ListMeta `json:"metadata"`
	Items       []IotDevice `json:"items"`
}

func (iotDevice *IotDevice) GetObjectKind() schema.ObjectKind {
	return &iotDevice.TypeMeta
}

func (iotDevice *IotDevice) GetObjectMeta() *v1.ListMeta {
	return &iotDevice.Metadata
}

func (iotDeviceList *IotDeviceList) GetObjectKind() schema.ObjectKind {
	return &iotDeviceList.TypeMeta
}

func (iotDeviceList *IotDeviceList) GetListMeta() v1.List {
	return &iotDeviceList.Metadata
}
