package kubernetes

import (
	"log"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/fields"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
)

func CreateIotPods(ds types.IotDaemonSet, dynamicClient *dynamic.Client,
	restClient *rest.RESTClient) error {
	var pods []types.IotPod
	devices, err := GetDaemonSetDevices(ds, dynamicClient, restClient)

	if err != nil {
		return err
	}

	// TODO check if pods don't exist already!

	for _, device := range devices {
		pod := types.IotPod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "IotPod",
				APIVersion: ds.APIVersion,
			},
			Metadata: v1.ObjectMeta{
				Name:      ds.Metadata.Name + "-" + string(common.NewUUID()),
				Namespace: ds.Metadata.Namespace,
				Labels: map[string]string{
					api.CreatedByAnnotation: ds.Metadata.SelfLink,
					types.DeviceSelector:    device.Metadata.Name,
				},
			},
			Spec: ds.Spec.Template.Spec,
		}
		pods = append(pods, pod)
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

// TODO Add function to retrieve related devices. Devices for pod can be discovered using
// "deviceSelector" label from pod (it's copied from daemon set during pod creation).
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

// TODO Add function to retrieve related daemon sets. Daemon sets can be discovered using
// "createdBy" label from pod.
func GetPodDaemonSet(restClient *rest.RESTClient, pod types.IotPod) (types.IotDaemonSet, error) {
	var ds types.IotDaemonSet

	fieldSelector, err := fields.ParseSelector("metadata.selfLink=" + pod.Metadata.Labels[api.CreatedByAnnotation])

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
