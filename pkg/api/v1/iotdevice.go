package v1

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/client-go/pkg/api/v1"
)

const (
	IotDeviceKind = "IotDevice"
	IotDeviceType = "iotdevices"
)

type IotDevice struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`
	Spec            v1.NodeSpec       `json:"spec"`
	Status          v1.NodeStatus     `json:"status"`
}

type IotDeviceList struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ListMeta `json:"metadata"`
	Items           []IotDevice     `json:"items"`
}

func (iotDevice *IotDevice) GetObjectKind() schema.ObjectKind {
	return &iotDevice.TypeMeta
}

func (iotDevice *IotDevice) GetObjectMeta() *metav1.ObjectMeta {
	return &iotDevice.Metadata
}

func (iotDeviceList *IotDeviceList) GetObjectKind() schema.ObjectKind {
	return &iotDeviceList.TypeMeta
}

func (iotDeviceList *IotDeviceList) GetListMeta() metav1.List {
	return &iotDeviceList.Metadata
}

func (this *IotDevice) UnmarshalJSON(data []byte) error {
	aux := fakeIotDevice{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	result := aux.toIotDevice()

	this.Status = result.Status
	this.Metadata = result.Metadata
	this.TypeMeta = result.TypeMeta
	this.Spec = result.Spec

	return nil
}

type fakeContainerImage struct {
	Names     []string `json:"names"`
	SizeBytes float64  `json:"sizeBytes,omitempty"`
}

func (this *fakeContainerImage) toContainerImage() v1.ContainerImage {
	return v1.ContainerImage{
		Names:     this.Names,
		SizeBytes: int64(this.SizeBytes),
	}
}

type fakeIotDevice struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`
	Spec            v1.NodeSpec       `json:"spec"`
	Status          struct {
		Capacity        v1.ResourceList        `json:"capacity,omitempty"`
		Allocatable     v1.ResourceList        `json:"allocatable,omitempty"`
		Phase           v1.NodePhase           `json:"phase,omitempty"`
		Conditions      []v1.NodeCondition     `json:"conditions,omitempty"`
		Addresses       []v1.NodeAddress       `json:"addresses,omitempty"`
		DaemonEndpoints v1.NodeDaemonEndpoints `json:"daemonEndpoints,omitempty"`
		NodeInfo        v1.NodeSystemInfo      `json:"nodeInfo,omitempty"`
		Images          []fakeContainerImage   `json:"images,omitempty"`
		VolumesInUse    []v1.UniqueVolumeName  `json:"volumesInUse,omitempty"`
		VolumesAttached []v1.AttachedVolume    `json:"volumesAttached,omitempty"`
	} `json:"status"`
}

func (this *fakeIotDevice) toIotDevice() *IotDevice {
	result := &IotDevice{
		TypeMeta: this.TypeMeta,
		Metadata: this.Metadata,
		Spec:     this.Spec,
		Status: v1.NodeStatus{
			Capacity:        this.Status.Capacity,
			Allocatable:     this.Status.Allocatable,
			Phase:           this.Status.Phase,
			Conditions:      this.Status.Conditions,
			Addresses:       this.Status.Addresses,
			DaemonEndpoints: this.Status.DaemonEndpoints,
			NodeInfo:        this.Status.NodeInfo,
			VolumesInUse:    this.Status.VolumesInUse,
			VolumesAttached: this.Status.VolumesAttached,
		},
	}

	for _, image := range this.Status.Images {
		result.Status.Images = append(result.Status.Images, image.toContainerImage())
	}

	return result
}
