package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/pkg/api/v1"
)

const IotPodType = "iotpods"

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
