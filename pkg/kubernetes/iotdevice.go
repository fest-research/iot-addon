package kubernetes

import (
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"fmt"
	"log"
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

func GetDaemonSetsForAllDevices(restClient *rest.RESTClient, namespace string) ([]types.IotDaemonSet, error) {
	return getGetDaemonSetsByDeviceSelector(restClient, namespace, types.DevicesAll)
}

func GetDaemonSetsForDevice(restClient *rest.RESTClient, namespace string, deviceName string) ([]types.IotDaemonSet, error) {
	return getGetDaemonSetsByDeviceSelector(restClient, namespace, deviceName)
}

func GetDeviceSelectedPods(device types.IotDevice, dynamicClient *dynamic.Client,
restClient *rest.RESTClient) ([]types.IotPod, error) {

	var dsList []types.IotDaemonSet
	var resultList []types.IotPod

	dsForAll, err := GetDaemonSetsForAllDevices(restClient, device.Metadata.Namespace)
	if err != nil {
		return nil, err
	}
	dsList = append(dsList, dsForAll...)

	dsForDevice, err := GetDaemonSetsForDevice(restClient, device.Metadata.Namespace, device.Metadata.Name)
	if err != nil {
		return nil, err
	}
	dsList = append(dsList, dsForDevice...)

	for _, item := range dsList {
		daemonSetPods, err := GetIotPods(dynamicClient, device.Metadata.Namespace, item.Metadata.SelfLink, device.Metadata.Name)
		if err != nil {
			log.Panic(fmt.Sprintf("GetDaemonSetPods error %v", err))
			continue
		}
		resultList = append(resultList, daemonSetPods...)
	}

	return resultList, nil
}

func getGetDaemonSetsByDeviceSelector(restClient *rest.RESTClient, namespace string, selectorValue string) ([]types.IotDaemonSet, error) {
	var dsList types.IotDaemonSetList
	err := restClient.Get().
		Resource(types.IotDaemonSetType).
		Namespace(namespace).LabelsSelectorParam(labels.Set{types.DeviceSelector: selectorValue}.AsSelector()).Do().
		Into(&dsList)

	if (err != nil) {
		return nil, err
	}

	return dsList.Items, nil
}

// TODO Add function to retrieve related pods. Pods for device can be discovered using
// "deviceSelector" label from pod (it's copied from daemon set during pod creation).

// TODO Add function to retrieve related daemon sets. Daemon sets can be discovered using
// "deviceSelector" label from daemon set.
