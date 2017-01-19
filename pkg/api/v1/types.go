package v1

import (
	"k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type IotDevice struct {
	metav1.TypeMeta `json:",inline"`
	Metadata    metav1.ListMeta `json:"metadata"`
}

type IotDeviceList struct {
	metav1.TypeMeta `json:",inline"`
	Metadata    metav1.ListMeta `json:"metadata"`
	Items       []IotDevice `json:"items"`
}

func (iotDevice *IotDevice) GetObjectKind() schema.ObjectKind {
	return &iotDevice.TypeMeta
}

func (iotDevice *IotDevice) GetObjectMeta() *metav1.ListMeta {
	return &iotDevice.Metadata
}

func (iotDeviceList *IotDeviceList) GetObjectKind() schema.ObjectKind {
	return &iotDeviceList.TypeMeta
}

func (iotDeviceList *IotDeviceList) GetListMeta() metav1.List {
	return &iotDeviceList.Metadata
}

type IotDaemonSet struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	// +optional
	Metadata v1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired behavior of this daemon set.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status
	// +optional
	Spec v1beta1.DaemonSetSpec `json:"spec,omitempty"`

	// Status is the current status of this daemon set. This data may be
	// out of date by some window of time.
	// Populated by the system.
	// Read-only.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status
	// +optional
	Status v1beta1.DaemonSetStatus `json:"status,omitempty"`
}

type IotDaemonSetList struct {
	metav1.TypeMeta `json:",inline"`

	Metadata metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of daemon sets.
	Items []IotDaemonSet `json:"items"`
}


func (iotDaemonSet *IotDaemonSet) GetObjectKind() schema.ObjectKind {
	return &iotDaemonSet.TypeMeta
}

func (iotDaemonSet *IotDaemonSet) GetObjectMeta() *v1.ObjectMeta {
	return &iotDaemonSet.Metadata
}

func (iotDaemonSetList *IotDaemonSetList) GetObjectKind() schema.ObjectKind {
	return &iotDaemonSetList.TypeMeta
}

func (iotDaemonSetList *IotDaemonSetList) GetListMeta() metav1.List {
	return &iotDaemonSetList.Metadata
}


// IotPod is a collection of containers that can run on a IotDevice.
type IotPod struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	// +optional
	Metadata v1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the pod.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status
	// +optional
	Spec v1.PodSpec `json:"spec,omitempty"`

	// Most recently observed status of the pod.
	// This data may not be up to date.
	// Populated by the system.
	// Read-only.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#spec-and-status
	// +optional
	Status v1.PodStatus `json:"status,omitempty"`
}

// IotPodList is a list of IotPods.
type IotPodList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#types-kinds
	// +optional
	Metadata metav1.ListMeta `json:"metadata,omitempty"`

	// List of iotPods.
	Items []IotPod `json:"items"`
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