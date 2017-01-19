package kubernetes

import (
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

// GetAllDevices returns all IotDevices from selected namespace.
func GetAllDevices(dynamicClient *dynamic.Client, namespace string) ([]types.IotDevice, error) {
	devices, err := dynamicClient.Resource(&metav1.APIResource{
		Name:       types.IotDeviceType,
		Namespaced: namespace != api.NamespaceNone,
	}, namespace).List(&v1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return devices.(*types.IotDeviceList).Items, nil
}

// GetAllDevices returns IotDevice with selected name from selected namespace.
func GetDevice(restClient *rest.RESTClient, name, namespace string) types.IotDevice {
	var device types.IotDevice

	restClient.Get().
		Resource(types.IotDeviceType).
		Namespace(namespace).
		Name(name).
		Do().
		Into(&device)

	return device
}

// TODO Add function to retrieve related pods. Pods for device can be discovered using
// "deviceSelector" label from pod (it's copied from daemon set during pod creation).

// TODO Add function to retrieve related daemon sets. Daemon sets can be discovered using
// "deviceSelector" label from daemon set.
