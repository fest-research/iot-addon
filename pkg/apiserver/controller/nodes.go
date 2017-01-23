package controller

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/emicklei/go-restful/log"
	"github.com/fest-research/iot-addon/pkg/api/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/watch"
	kubeapi "k8s.io/client-go/pkg/api/v1"
)

type NodeController struct{}

func (this *NodeController) Transform(in interface{}) (interface{}, error) {
	log.Print("NodeController - Transform()")

	switch in.(type) {
	case watch.Event:
		event := in.(watch.Event)
		return this.transformWatchEvent(event), nil
	case *v1.IotDeviceList:
		iotDeviceList := in.(*v1.IotDeviceList)
		return this.toNodeList(iotDeviceList), nil
	case *v1.IotDevice:
		iotDevice := in.(*v1.IotDevice)
		return this.toNode(iotDevice), nil
	case *kubeapi.Node:
		node := in.(*kubeapi.Node)
		return this.toUnstructured(node)
	case *unstructured.Unstructured:
		unstructured := in.(*unstructured.Unstructured)
		return this.toBytes(unstructured)
	default:
		return nil, fmt.Errorf("Not supported type: %s", reflect.TypeOf(in))
	}
}

func (this *NodeController) transformWatchEvent(event watch.Event) watch.Event {
	iotDevice := event.Object.(*v1.IotDevice)
	event.Object = this.toNode(iotDevice)
	return event
}

func (this *NodeController) toNodeList(iotDeviceList *v1.IotDeviceList) kubeapi.NodeList {
	nodeList := kubeapi.NodeList{}

	nodeList.Kind = "NodeList"
	nodeList.APIVersion = "v1"
	nodeList.Items = make([]kubeapi.Node, 0)

	for _, iotDevice := range iotDeviceList.Items {
		node := this.toNode(&iotDevice)
		nodeList.Items = append(nodeList.Items, *node)
	}

	return nodeList
}

func (this *NodeController) toNode(iotDevice *v1.IotDevice) *kubeapi.Node {
	node := &kubeapi.Node{}

	// TODO: subject to revision
	node.Kind = "Node"
	node.APIVersion = "v1"
	node.Spec = iotDevice.Spec
	node.Status = iotDevice.Status
	node.ObjectMeta = iotDevice.Metadata

	return node
}

func (this *NodeController) toIotDevice(node *kubeapi.Node) *v1.IotDevice {
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
func (this *NodeController) toUnstructured(node *kubeapi.Node) (*unstructured.Unstructured, error) {
	result := &unstructured.Unstructured{}
	iotDevice := this.toIotDevice(node)

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
func (this *NodeController) toBytes(unstructured *unstructured.Unstructured) ([]byte, error) {
	marshalledIotDevice, err := unstructured.MarshalJSON()
	if err != nil {
		return nil, err
	}

	iotDevice := &v1.IotDevice{}
	err = json.Unmarshal(marshalledIotDevice, iotDevice)
	if err != nil {
		return nil, err
	}

	node := this.toNode(iotDevice)
	marshalledNode, err := json.Marshal(node)
	if err != nil {
		return nil, err
	}

	return marshalledNode, nil
}

func NewNodeController() *NodeController {
	return &NodeController{}
}
