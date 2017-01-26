package controller

import (
	"github.com/fest-research/iot-addon/pkg/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"

	kubeapi "k8s.io/client-go/pkg/api/v1"
)

type INodeController interface {
	// TransformWatchEvent implements WatchEventController.
	TransformWatchEvent(event watch.Event) watch.Event
	ToNodeList(*v1.IotDeviceList) *kubeapi.NodeList
	ToNode(*v1.IotDevice) *kubeapi.Node
	ToIotDevice(*kubeapi.Node) *v1.IotDevice
	ToUnstructured(*kubeapi.Node) (*unstructured.Unstructured, error)
	ToBytes(*unstructured.Unstructured) ([]byte, error)
}

type nodeController struct{}

func (this nodeController) TransformWatchEvent(event watch.Event) watch.Event {
	iotDevice := event.Object.(*v1.IotDevice)
	event.Object = this.ToNode(iotDevice)
	return event
}

func (this nodeController) ToNodeList(iotDeviceList *v1.IotDeviceList) *kubeapi.NodeList {
	nodeList := &kubeapi.NodeList{}

	nodeList.TypeMeta = this.getTypeMeta(v1.NodeListKind)
	nodeList.Items = make([]kubeapi.Node, 0)

	for _, iotDevice := range iotDeviceList.Items {
		node := this.ToNode(&iotDevice)
		nodeList.Items = append(nodeList.Items, *node)
	}

	return nodeList
}

func (this nodeController) ToNode(iotDevice *v1.IotDevice) *kubeapi.Node {
	node := &kubeapi.Node{}

	// TODO: subject to revision
	node.TypeMeta = this.getTypeMeta(v1.NodeKind)

	node.Spec = iotDevice.Spec
	node.Status = iotDevice.Status
	node.ObjectMeta = iotDevice.Metadata

	node.ObjectMeta.Namespace = ""

	return node
}

func (this nodeController) ToIotDevice(node *kubeapi.Node) *v1.IotDevice {
	iotDevice := &v1.IotDevice{}

	// TODO: subject to revision
	iotDevice.TypeMeta = this.getIotTypeMeta()

	iotDevice.Metadata = node.ObjectMeta
	iotDevice.Status = node.Status
	iotDevice.Spec = node.Spec

	return iotDevice
}

// Converts node to unstructured iot device
func (this nodeController) ToUnstructured(node *kubeapi.Node) (*unstructured.Unstructured, error) {
	result := &unstructured.Unstructured{}
	iotDevice := this.ToIotDevice(node)

	marshalledIotDevice, err := json.Marshal(iotDevice)
	if err != nil {
		return nil, err
	}

	err = result.UnmarshalJSON(marshalledIotDevice)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Converts unstructured iot device to node json bytes array
func (this nodeController) ToBytes(unstructured *unstructured.Unstructured) ([]byte, error) {
	marshalledIotDevice, err := unstructured.MarshalJSON()
	if err != nil {
		return nil, err
	}

	iotDevice := &v1.IotDevice{}
	err = json.Unmarshal(marshalledIotDevice, iotDevice)
	if err != nil {
		return nil, err
	}

	node := this.ToNode(iotDevice)
	marshalledNode, err := json.Marshal(node)
	if err != nil {
		return nil, err
	}

	return marshalledNode, nil
}

func (this nodeController) getIotTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: v1.IotAPIVersion,
		Kind:       v1.IotDeviceKind,
	}
}

func (this nodeController) getTypeMeta(kind v1.ResourceKind) metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: v1.APIVersion,
		Kind:       string(kind),
	}
}

func NewNodeController() INodeController {
	return &nodeController{}
}
