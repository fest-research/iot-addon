package kubernetes

import (
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// GetDaemonSetSelectedDevices returns all IotDevices from selected namespace, that are specified
// in IotDaemonSet with deviceSelector (IotDevice name or 'all').
func GetDaemonSetDevices(ds types.IotDaemonSet, dynamicClient *dynamic.Client,
	restClient *rest.RESTClient) ([]types.IotDevice, error) {
	deviceSelector := ds.Metadata.Labels[types.DeviceSelector]
	if deviceSelector == types.DevicesAll {
		return GetAllDevices(dynamicClient, ds.Metadata.Namespace)
	} else {

		result := []types.IotDevice{}
		device, err := GetDevice(restClient, deviceSelector, ds.Metadata.Namespace)
		if err != nil {
			return result, nil
		}
		result = append(result, device)
		return result, nil
	}
}

func GetDaemonSetPods(restClient *rest.RESTClient, ds types.IotDaemonSet) ([]types.IotPod, error) {
	var podList types.IotPodList
	err := restClient.Get().
		Resource(types.IotPodType).
		Namespace(ds.Metadata.Namespace).
		LabelsSelectorParam(labels.Set{
			types.CreatedBy: types.IotDaemonSetType + "." + ds.Metadata.Name,
		}.AsSelector()).
		Do().
		Into(&podList)
	return podList.Items, err
}

func DaemonSetToPod(ds types.IotDaemonSet) types.IotPod {
	deviceSelector, ok := ds.Metadata.Labels[types.DeviceSelector]

	if !ok {
		deviceSelector = ""
	}

	return types.IotPod{
		TypeMeta: metav1.TypeMeta{
			Kind:       types.IotPodKind,
			APIVersion: ds.APIVersion,
		},
		Metadata: metav1.ObjectMeta{
			Name:      ds.Metadata.Name,
			Namespace: ds.Metadata.Namespace,
			Labels: map[string]string{
				types.CreatedBy:      types.IotDaemonSetType + "." + ds.Metadata.Name,
				types.DeviceSelector: deviceSelector,
			},
		},
		Spec: ds.Spec.Template.Spec,
	}
}
