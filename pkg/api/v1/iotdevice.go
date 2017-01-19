package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/pkg/api/v1"
)

// TODO Add function to retrieve related pods. Pods for device can be discovered using
// "deviceSelector" label from pod (it's copied from daemon set during pod creation).

// TODO Add function to retrieve related daemon sets. Daemon sets can be discovered using
// "deviceSelector" label from daemon set.

type IotDevice struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        v1.ObjectMeta `json:"metadata"`
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
