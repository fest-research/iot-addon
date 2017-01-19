package kubernetes

import (
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

func GetAllDevices(dynamicClient *dynamic.Client, namespace string) ([]types.IotDevice, error) {
	devices, err := dynamicClient.Resource(&metav1.APIResource{
		Name:       "iotdevices",
		Namespaced: namespace != api.NamespaceNone,
	}, namespace).List(&v1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return devices.(*types.IotDeviceList).Items, nil
}

func GetDevice(restClient *rest.RESTClient, name, namespace string) types.IotDevice {
	var device types.IotDevice

	restClient.Get().
		Resource("iotdevices").
		Namespace(namespace).
		Name(name).
		Do().
		Into(&device)

	return device
}
