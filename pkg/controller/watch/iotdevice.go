package watch

import (
	"log"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"time"
)

type IotDeviceWatcher struct {
	dynamicClient *dynamic.Client
	restClient    *rest.RESTClient
}

var iotDeviceResource = metav1.APIResource{
	Name:       types.IotDeviceType,
	Namespaced: true,
}

func NewIotDeviceWatcher(dynamicClient *dynamic.Client, restClient *rest.RESTClient) IotDeviceWatcher {
	return IotDeviceWatcher{dynamicClient: dynamicClient, restClient: restClient}
}

func (w IotDeviceWatcher) Watch() {

	var watcher watch.Interface = nil
	var err error = nil
	ticker := time.NewTicker(time.Second * 4)
	defer ticker.Stop()

	for ok := true; ok; ok = (watcher == nil) {
		select {
		case <-ticker.C:
			watcher, err = w.dynamicClient.
				Resource(&iotDeviceResource, api.NamespaceAll).
				Watch(&api.ListOptions{})
			if err != nil {
				log.Println(err.Error())
			} else {
				ticker.Stop()
			}
			break
		}
	}

	defer watcher.Stop()
	log.Printf("Watcher for %s created \n", types.IotDeviceType)

	for {
		e, ok := <-watcher.ResultChan()

		if !ok {
			panic("IotDevices ended early?")
		}

		modificationMap := map[string]bool{}
		iotDevice, _ := e.Object.(*types.IotDevice)

		if e.Type == watch.Added {
			log.Printf("Device added %s\n", iotDevice.Metadata.Name)
			unschedulable := kubernetes.GetUnschedulableLabelFromDevice(*iotDevice)
			modificationMap[iotDevice.Metadata.Name] = unschedulable
			err := w.addModifyDeviceHandler(*iotDevice)
			if err != nil {
				log.Printf("Error [addModifyDeviceHandler] %s", err.Error())
			}
		} else if e.Type == watch.Modified {
			unschedulable := kubernetes.GetUnschedulableLabelFromDevice(*iotDevice)
			prevUnschedulable := modificationMap[iotDevice.Metadata.Name]
			if unschedulable != prevUnschedulable {
				log.Printf("Device  modified %s\n", iotDevice.Metadata.Name)
				err := w.addModifyDeviceHandler(*iotDevice)
				if err != nil {
					log.Printf("Error [addModifyDeviceHandler] %s", err.Error())
				}
				modificationMap[iotDevice.Metadata.Name] = unschedulable
			}

		} else if e.Type == watch.Error {
			log.Println("Error")
			break
		}
	}
}

func (w IotDeviceWatcher) addModifyDeviceHandler(iotDevice types.IotDevice) error {

	unschedulable := kubernetes.GetUnschedulableLabelFromDevice(iotDevice)
	deviceName := iotDevice.Metadata.Name

	if unschedulable {
		log.Printf("[addModifyDeviceHandler] Delete pods for unschedulable device %s", deviceName)
		pods, err := kubernetes.GetDevicePods(w.restClient, iotDevice)
		if err != nil {
			return err
		}

		for _, pod := range pods {
			err := w.deletePod(pod)
			if err != nil {
				return err
			}
		}

	} else {
		daemonSets, _ := kubernetes.GetDeviceDaemonSets(w.restClient, iotDevice)
		for _, ds := range daemonSets {

			if !kubernetes.IsPodCreated(w.restClient, ds, iotDevice) {
				log.Printf("[addModifyDeviceHandler] Create new pod %s ", ds.Metadata.Name)
				err := kubernetes.CreateDaemonSetPod(ds, iotDevice, w.restClient)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (w IotDeviceWatcher) deletePod(pod types.IotPod) error {

	return w.restClient.Delete().
		Namespace(pod.Metadata.Namespace).
		Resource(types.IotPodType).
		Name(pod.Metadata.Name).
		Body(&v1.DeleteOptions{}).
		Do().
		Error()

}

func createTypeMeta(apiVersion string) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       types.IotPodKind,
		APIVersion: apiVersion,
	}
}
