package kubernetes

import (
	"log"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/rest"
)

// CreateDaemonSetPod created ds pod on specific device
func CreateDaemonSetPod(ds types.IotDaemonSet, device types.IotDevice, restClient *rest.RESTClient) error {
	labels := map[string]string{
		types.CreatedBy:      types.IotDaemonSetType + "." + ds.Metadata.Name,
		types.DeviceSelector: device.Metadata.Name,
	}
	common.MapCopy(labels, ds.Spec.Template.ObjectMeta.Labels)
	return restClient.Post().
		Namespace(ds.Metadata.Namespace).
		Resource(types.IotPodType).
		Body(&types.IotPod{
			TypeMeta: metav1.TypeMeta{
				Kind:       types.IotPodKind,
				APIVersion: ds.APIVersion,
			},
			Metadata: v1.ObjectMeta{
				Name:      ds.Metadata.Name + "-" + string(common.NewUUID()), // TODO use template val
				Namespace: ds.Metadata.Namespace,
				Labels:    labels,
			},
			Spec: ds.Spec.Template.Spec,
		}).Do().Error()
}

// isPodCorrectlyScheduled checks if pod is correctly scheduled.
func IsPodCorrectlyScheduled(ds types.IotDaemonSet, pod types.IotPod) bool {
	if ds.Metadata.Labels[types.DeviceSelector] == types.DevicesAll {
		return true
	} else {
		return ds.Metadata.Labels[types.DeviceSelector] == pod.Metadata.Labels[types.DeviceSelector] &&
			ds.Metadata.Namespace == pod.Metadata.Namespace
	}
}

// getDevicesMissingPods filters daemon set destined devices and returns devices without any existing pods.
func GetDevicesMissingPods(dsDestinedDevices []types.IotDevice, existingPods []types.IotPod) []types.IotDevice {
	var devicesMissingPod []types.IotDevice
	for _, device := range dsDestinedDevices {
		unschedulable := GetUnschedulableLabelFromDevice(device)
		if !unschedulable {
			// Assume that device has missing pod.
			isPodMissing := true
			for _, pod := range existingPods {
				// Check if it really has iterating through all pods and checking them.
				if pod.Metadata.Labels[types.DeviceSelector] == device.Metadata.Name {
					isPodMissing = false
					break
				}
			}

			// If device has missing pod, then it has to be added to output array.
			if isPodMissing {
				devicesMissingPod = append(devicesMissingPod, device)
			}
		}
	}
	return devicesMissingPod

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

func DeletePod(restClient *rest.RESTClient, pod types.IotPod) {
	restClient.Delete().Resource(types.IotPodType).Namespace(pod.Metadata.Namespace).Name(pod.Metadata.Name).Do()
}

// TODO Update name, namespace and labels?
// TODO Save updated pod.
func UpdatePod(restClient *rest.RESTClient, pod types.IotPod, template v1.PodTemplateSpec) {
	newPod := types.IotPod{}
	pod.Spec = template.Spec

	labels := map[string]string{
		types.CreatedBy:      pod.Metadata.Labels[types.CreatedBy],
		types.DeviceSelector: pod.Metadata.Labels[types.DeviceSelector],
	}
	common.MapCopy(labels, template.ObjectMeta.Labels)

	pod.Metadata.Labels = labels

	err := restClient.Put().
		Namespace(pod.Metadata.Namespace).
		Resource(types.IotPodType).
		Name(pod.Metadata.Name).
		Body(&pod).
		Do().
		Into(&newPod)
	if err != nil {
		log.Printf("Error. Can not update IotPod %s", pod.Metadata.Name)
	}

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
