package kubernetes

import (
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"log"
)

func CreateIotPods(ds types.IotDaemonSet, dynamicClient *dynamic.Client,
	restClient *rest.RESTClient) error {
	var pods []types.IotPod
	devices, err := GetDaemonSetSelectedDevices(ds, dynamicClient, restClient)

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
				Annotations: map[string]string{
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
			Resource("iotpods").
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

// TODO Add function to retrieve related daemon sets. Daemon sets can be discovered using
// "createdBy" label from pod.
