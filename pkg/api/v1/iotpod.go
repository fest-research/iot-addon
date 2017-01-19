package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/pkg/api/v1"
)

// TODO Add function to retrieve related devices. Devices for pod can be discovered using
// "deviceSelector" label from pod (it's copied from daemon set during pod creation).

// TODO Add function to retrieve related daemon sets. Daemon sets can be discovered using
// "createdBy" label from pod.

type IotPod struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        v1.ObjectMeta `json:"metadata,omitempty"`
	Spec            v1.PodSpec    `json:"spec,omitempty"`
	Status          v1.PodStatus  `json:"status,omitempty"`
}

type IotPodList struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IotPod        `json:"items"`
}

func (iotPod *IotPod) GetObjectKind() schema.ObjectKind {
	return &iotPod.TypeMeta
}

func (iotPod *IotPod) GetObjectMeta() *v1.ObjectMeta {
	return &iotPod.Metadata
}

func (iotPodList *IotPodList) GetObjectKind() schema.ObjectKind {
	return &iotPodList.TypeMeta
}

func (iotPodList *IotPodList) GetListMeta() metav1.List {
	return &iotPodList.Metadata
}
