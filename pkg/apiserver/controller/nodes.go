package controller

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/emicklei/go-restful/log"
	"github.com/fest-research/iot-addon/pkg/api/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kubeapi "k8s.io/client-go/pkg/api/v1"
)

type NodeController struct{}

func (this *NodeController) Transform(in interface{}) (interface{}, error) {
	log.Print("NodeController - Transform()")

	switch in.(type) {
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

func (this *NodeController) toNode(iotDevice *v1.IotDevice) *kubeapi.Node {
	node := &kubeapi.Node{}

	// TODO: transform the iotDevice into the node

	return node
}

func (this *NodeController) toIotDevice(node *kubeapi.Node) *v1.IotDevice {
	iotDevice := &v1.IotDevice{}

	// TODO: decide which other fields to propagate to the IotDevice
	iotDevice.Kind = "IotDevice"
	iotDevice.APIVersion = "fujitsu.com/v1"
	iotDevice.Metadata = node.ObjectMeta
	iotDevice.Status = node.Status

	return iotDevice
}

// Converts pod to unstructured iot pod
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

// Converts unstructured iot pod to pod json bytes array
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
