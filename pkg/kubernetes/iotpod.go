package kubernetes

import (
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/rest"
	"log"
)

// CreatePods creates IotPods for IotDaemonSet if they already don't exist.
func CreatePods(ds types.IotDaemonSet, dynamicClient *dynamic.Client, restClient *rest.RESTClient) error {
	var pods []types.IotPod

	devices, err := GetDaemonSetDevices(ds, dynamicClient, restClient)

	if err != nil {
		return err
	}

	for _, device := range devices {
		if !IsPodCreated(restClient, ds, device) {
			pods = append(pods, types.IotPod{
				TypeMeta: metav1.TypeMeta{
					Kind:       "IotPod",
					APIVersion: ds.APIVersion,
				},
				Metadata: v1.ObjectMeta{
					Name:      ds.Metadata.Name + "-" + string(common.NewUUID()),
					Namespace: ds.Metadata.Namespace,
					Labels: map[string]string{
						types.CreatedBy:      types.IotDaemonSetType + "." + ds.Metadata.Name,
						types.DeviceSelector: device.Metadata.Name,
					},
				},
				Spec: ds.Spec.Template.Spec,
			})
		}
	}

	for _, pod := range pods {
		newPod := types.IotPod{}

		err = restClient.Post().
			Namespace(ds.Metadata.Namespace).
			Resource(types.IotPodType).
			Body(&pod).
			Do().
			Into(&newPod)

		log.Printf("Created new pod %s for %s daemon set",
			newPod.Metadata.Name,
			ds.Metadata.Name)
	}

	return nil
}

// IsPodCreated checks if there is any IotPod created for IotDaemonSet on IotDevice.
// TODO Check should check if IotPods are the same if they don't exist and update them then.
func IsPodCreated(restClient *rest.RESTClient, ds types.IotDaemonSet, device types.IotDevice) bool {
	var podList types.IotPodList

	// Ignores error.
	restClient.Get().
		Resource(types.IotPodType).
		Namespace(ds.Metadata.Namespace).
		LabelsSelectorParam(labels.Set{
			types.CreatedBy:      types.IotDaemonSetType + "." + ds.Metadata.Name,
			types.DeviceSelector: device.Metadata.Name,
		}.AsSelector()).
		Do().
		Into(&podList)

	return len(podList.Items) > 0
}

// GetPodDevice returns IotDevice where IotPod is deployed. Method uses "deviceSelector" label from IotPod.
func GetPodDevice(restClient *rest.RESTClient, pod types.IotPod) (types.IotDevice, error) {
	var device types.IotDevice

	fieldSelector, err := fields.ParseSelector("metadata.name=" + pod.Metadata.Labels[types.DeviceSelector])

	if err != nil {
		return device, err
	}

	err = restClient.Get().
		Resource(types.IotDeviceType).
		Namespace(device.Metadata.Namespace).
		FieldsSelectorParam(fieldSelector).
		Do().
		Into(&device)

	return device, err
}

// GetPodDaemonSet returns IotDaemonSet which created IotPod. Method uses "createdBy" label from IotPod.
func GetPodDaemonSet(restClient *rest.RESTClient, pod types.IotPod) (types.IotDaemonSet, error) {
	var ds types.IotDaemonSet

	fieldSelector, err := fields.ParseSelector("metadata.selfLink=" + pod.Metadata.Labels[types.CreatedBy])

	if err != nil {
		return ds, err
	}

	err = restClient.Get().
		Resource(types.IotDaemonSetType).
		Namespace(ds.Metadata.Namespace).
		FieldsSelectorParam(fieldSelector).
		Do().
		Into(&ds)

	return ds, err
}
