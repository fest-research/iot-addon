package kubernetes

import (
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
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

// TODO Add function to retrieve related pods. Pods for daemon set can be discovered using
// "createdBy" annotation from pod.
