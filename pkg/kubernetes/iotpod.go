package kubernetes

import (
	"log"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/rest"
)

// CreateDaemonSetPods creates pods for daemon set if they already don't exist.
// TODO Split into template filling and object creation methods.
// TODO Create "kind" fields for custom types.
func CreateDaemonSetPods(ds types.IotDaemonSet, dynamicClient *dynamic.Client, restClient *rest.RESTClient) error {
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

		log.Printf("Created new pod %s for %s daemon set", newPod.Metadata.Name, ds.Metadata.Name)
	}

	return nil
}

// TODO interface for pods

// UpdateDaemonSetPods updates all daemon set pods after it was modified. Pods have to be correctly scheduled and have
// up-to-date specs.
func UpdateDaemonSetPods(restClient *rest.RESTClient, dynamicClient *dynamic.Client, ds types.IotDaemonSet) error {
	log.Printf("Updating pods created by %s %s\n", ds.Metadata.Name, ds.TypeMeta.Kind)

	// Making sure, that daemon set is deployed on currently selected devices.
	// Getting all existing pods created by daemon set.
	existingPods, err := GetDaemonSetPods(restClient, ds)
	if err != nil {
		return err
	}

	// Getting list of devices where daemon set should be deployed.
	destinedDevices, err := GetDaemonSetDevices(ds, dynamicClient, restClient)

	// Updating existing pods.
	for _, existingPod := range existingPods {
		if !shouldBePodScheduled(destinedDevices, ds, existingPod) {
			// TODO remove pod
		} else {
			// TODO update pod spec
		}
	}

	// TODO add missing

	return nil
}

// shouldBePodScheduled checks if pod should be scheduled on device.
func shouldBePodScheduled(dsDestinedDevices []types.IotDevice, ds types.IotDaemonSet, pod types.IotPod) bool {
	if ds.Metadata.Labels[types.DeviceSelector] == types.DevicesAll {
		return true
	} else {
		for _, dsDestinedDevice := range dsDestinedDevices {
			if dsDestinedDevice.Metadata.Name == pod.Metadata.Labels[types.DeviceSelector] &&
				dsDestinedDevice.Metadata.Namespace == pod.Metadata.Namespace {
				return true
			}
		}
		return false
	}
}

func DeleteDaemonSetPods(restClient *rest.RESTClient, ds types.IotDaemonSet) error {
	log.Printf("Deleting pods created by %s %s\n", ds.Metadata.Name, ds.TypeMeta.Kind)
	return restClient.Delete().
		Resource(types.IotPodType).
		Namespace(ds.Metadata.Namespace).
		LabelsSelectorParam(labels.Set{
			types.CreatedBy: types.IotDaemonSetType + "." + ds.Metadata.Name,
		}.AsSelector()).
		Do().
		Error()
}

// IsPodCreated checks if there is any IotPod created for IotDaemonSet on IotDevice.
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
