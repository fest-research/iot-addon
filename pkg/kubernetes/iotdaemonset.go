package kubernetes

import (
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"log"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"fmt"

)



// GetDaemonSetSelectedDevices returns all IotDevices from selected namespace, that are specified
// in IotDaemonSet with deviceSelector (IotDevice name or 'all').
func GetDaemonSetSelectedDevices(ds types.IotDaemonSet, dynamicClient *dynamic.Client,
	restClient *rest.RESTClient) ([]types.IotDevice, error) {

	deviceSelector := ds.Metadata.Labels[types.DeviceSelector]
	if deviceSelector == types.DevicesAll {
		return GetAllDevices(dynamicClient, ds.Metadata.Namespace)
	} else {
		return []types.IotDevice{
			GetDevice(restClient, deviceSelector, ds.Metadata.Namespace),
		}, nil
	}
}

func GetDaemonSetSelectedPods(ds types.IotDaemonSet, dynamicClient *dynamic.Client,
restClient *rest.RESTClient) ([]types.IotPod, error) {

	var resultList []types.IotPod

	deviceList, err := GetDaemonSetSelectedDevices(ds, dynamicClient, restClient)
	if err != nil {
		return []types.IotPod{}, nil
	}

	for _, device := range deviceList {

		daemonSetPods, err := GetIotPods(dynamicClient, ds.Metadata.Namespace, ds.Metadata.SelfLink, device.Metadata.Name)
		if err != nil {
			log.Panic(fmt.Sprintf("GetDaemonSetPods error %v", err))
			continue
		}
		resultList = append(resultList, daemonSetPods...)
	}
	return resultList, nil
}
