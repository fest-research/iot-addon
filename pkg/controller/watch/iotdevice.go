package watch

import (
	"log"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	"github.com/fest-research/iot-addon/pkg/common"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
	"strconv"
)

var iotDeviceResource = metav1.APIResource{
	Name:       "iotdevices",
	Namespaced: true,
}

func WatchIotDevices(dynamicClient *dynamic.Client, restClient *rest.RESTClient) {
	watcher, err := dynamicClient.
	Resource(&iotDeviceResource, api.NamespaceAll).
		Watch(&api.ListOptions{})

	if err != nil {
		log.Println(err.Error())
	}

	defer watcher.Stop()

	for {
		e, ok := <-watcher.ResultChan()

		if !ok {
			panic("IotDevices ended early?")
		}

		iotDevice, _ := e.Object.(*types.IotDevice)

		if e.Type == watch.Added {
			log.Printf("Device added %s\n", iotDevice.Metadata.Name)
			err := addDeviceHandler(restClient, *iotDevice)
			if err != nil {
				log.Printf("Error [addDeviceHandler] %s", err.Error())
			}
		} else if e.Type == watch.Modified {
			log.Printf("Modified %s\n", iotDevice.Metadata.SelfLink)
		} else if e.Type == watch.Deleted {
			log.Printf("Deleted %s\n", iotDevice.Metadata.SelfLink)
		} else if e.Type == watch.Error {
			log.Println("Error")
			break
		}
	}
}

func addDeviceHandler(restClient *rest.RESTClient, iotDevice types.IotDevice) error {
	var unschedulable bool
	var err error
	daemonSets, _ := kubernetes.GetDeviceDaemonSets(restClient, iotDevice)

	unschedulableLabel, ok := iotDevice.Metadata.Labels[types.Unschedulable]

	if ok {
		unschedulable, err = strconv.ParseBool(unschedulableLabel)
		if err != nil {
			return err
		}
	}

	deviceName := iotDevice.Metadata.Name

	if unschedulable {
		log.Printf("[addDeviceHandler] Delete pods for unschedulable device %s", deviceName)
		pods, err := kubernetes.GetDevicePods(restClient, iotDevice)
		if err != nil {
			return err
		}

		for _, pod := range pods {
			err := deletePod(restClient, pod)
			if err != nil {
				return err
			}
		}

	} else {
		for _, ds := range daemonSets {

			if !kubernetes.IsPodCreated(restClient, ds, iotDevice) {
				log.Printf("[addDeviceHandler] Create new pod %s ", ds.Metadata.Name)
				err := createPod(restClient, ds, deviceName)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func createPod(restClient *rest.RESTClient, ds types.IotDaemonSet, deviceName string) error {

	newPod := types.IotPod{}
	name := ds.Metadata.Name

	pod := types.IotPod{
		TypeMeta: createTypeMeta(ds.APIVersion),
		Metadata: v1.ObjectMeta{
			Name:      name + "-" + string(common.NewUUID()),
			Namespace: ds.Metadata.Namespace,
			Labels: map[string]string{
				types.CreatedBy:      types.IotDaemonSetType + "." + name,
				types.DeviceSelector: deviceName,
			},
		},
		Spec: ds.Spec.Template.Spec,
	}

	return restClient.Post().
		Namespace(pod.Metadata.Namespace).
		Resource(types.IotPodType).
		Body(&pod).
		Do().
		Into(&newPod)

}

func deletePod(restClient *rest.RESTClient, pod types.IotPod) error {

	return restClient.Delete().
		Namespace(pod.Metadata.Namespace).
		Resource(types.IotPodType).
		Name(pod.Metadata.Name).
		Body(&v1.DeleteOptions{}).
		Do().
		Error()

}

func createTypeMeta(apiVersion string) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       "IotPod",
		APIVersion: apiVersion,
	}
}
