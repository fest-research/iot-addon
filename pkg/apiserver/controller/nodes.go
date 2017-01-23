package controller

import (
	"encoding/json"

	"github.com/fest-research/iot-addon/pkg/api/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

type NodeController struct{}

func (this NodeController) TransformWatchEvent(event watch.Event) watch.Event {
	iotDevice := event.Object.(*v1.IotDevice)
	event.Object = this.ToNode(iotDevice)
	return event
}

func (this NodeController) ToNodeList(iotDeviceList *v1.IotDeviceList) *kubeapi.NodeList {
	nodeList := &kubeapi.NodeList{}

	nodeList.Kind = "NodeList"
	nodeList.APIVersion = "v1"
	nodeList.Items = make([]kubeapi.Node, 0)

	for _, iotDevice := range iotDeviceList.Items {
		node := this.ToNode(&iotDevice)
		nodeList.Items = append(nodeList.Items, *node)
	}

	return nodeList
}

func (this NodeController) ToNode(iotDevice *v1.IotDevice) *kubeapi.Node {
	node := &kubeapi.Node{}

	// TODO: subject to revision
	node.Kind = "Node"
	node.APIVersion = "v1"
	node.Spec = iotDevice.Spec
	node.Status = iotDevice.Status
	node.ObjectMeta = iotDevice.Metadata

	return node
}

func (this NodeController) ToIotDevice(node *kubeapi.Node) *v1.IotDevice {
	iotDevice := &v1.IotDevice{}

	// TODO: subject to revision
	iotDevice.Kind = "IotDevice"
	iotDevice.APIVersion = "fujitsu.com/v1"
	iotDevice.Metadata = node.ObjectMeta
	iotDevice.Status = node.Status
	iotDevice.Spec = node.Spec

	return iotDevice
}

// Converts node to unstructured iot device
func (this NodeController) ToUnstructured(node *kubeapi.Node) (*unstructured.Unstructured, error) {
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
func (this NodeController) ToBytes(unstructured *unstructured.Unstructured) ([]byte, error) {
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

func NewNodeController() INodeController {
	return &NodeController{}
}
