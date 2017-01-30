package kubernetes

import (
	"strconv"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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
	return devices.(*types.IotDeviceList).Items, err
}

// GetAllDevices returns IotDevice with selected name from selected namespace.
func GetDevice(restClient *rest.RESTClient, name, namespace string) (types.IotDevice, error) {
	var device types.IotDevice
	err := restClient.Get().
		Resource(types.IotDeviceType).
		Namespace(namespace).
		Name(name).
		Do().
		Into(&device)

	if err != nil {
		return types.IotDevice{}, err
	}
	return device, nil
}

// TODO ?
func GetDeviceDaemonSets(restClient *rest.RESTClient, device types.IotDevice) ([]types.IotDaemonSet, error) {
	var dsList types.IotDaemonSetList
	var resList []types.IotDaemonSet

	err := restClient.Get().
		Resource(types.IotDaemonSetType).
		Namespace(device.Metadata.Namespace).
		LabelsSelectorParam(labels.Set{types.DeviceSelector: device.Metadata.Name}.
			AsSelector()).
		Do().
		Into(&dsList)

	if err != nil {
		return nil, err
	}
	resList = append(resList, dsList.Items...)

	err = restClient.Get().
		Resource(types.IotDaemonSetType).
		Namespace(device.Metadata.Namespace).
		LabelsSelectorParam(labels.Set{types.DeviceSelector: types.DevicesAll}.
			AsSelector()).
		Do().
		Into(&dsList)

	if err != nil {
		return nil, err
	}

	resList = append(resList, dsList.Items...)

	return resList, nil
}

func GetDevicePods(restClient *rest.RESTClient, device types.IotDevice) ([]types.IotPod, error) {
	var podList types.IotPodList
	err := restClient.Get().
		Resource(types.IotPodType).
		Namespace(device.Metadata.Namespace).
		LabelsSelectorParam(labels.Set{types.DeviceSelector: device.Metadata.Name}.
			AsSelector()).
		Do().
		Into(&podList)
	return podList.Items, err
}

func GetUnschedulableLabelFromDevice(iotDevice types.IotDevice) bool {

	unschedulableLabel, ok := iotDevice.Metadata.Labels[types.Unschedulable]

	if ok {
		unschedulable, err := strconv.ParseBool(unschedulableLabel)
		if err != nil {
			return false
		}
		return unschedulable
	}
	return false
}
