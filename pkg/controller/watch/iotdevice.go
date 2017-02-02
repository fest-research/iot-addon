package watch

import (
	"fmt"
	"log"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
)

type IotDeviceWatcher struct {
	dynamicClient *dynamic.Client
	restClient    *rest.RESTClient
	clientset     *client.Clientset
	iotDomain     string
}

func NewIotDeviceWatcher(dynamicClient *dynamic.Client, restClient *rest.RESTClient, clientset *client.Clientset,
	iotDomain string) IotDeviceWatcher {
	return IotDeviceWatcher{
		dynamicClient: dynamicClient,
		restClient:    restClient,
		clientset:     clientset,
		iotDomain:     iotDomain,
	}
}

func (w IotDeviceWatcher) Watch() {
	for {
		err := w.start()
		if err != nil {
			log.Printf("An error occured: %s", err.Error())
		}
	}
}

func (w IotDeviceWatcher) start() error {
	watcher, err := w.dynamicClient.Resource(&metav1.APIResource{
		Name:       types.IotDeviceType,
		Namespaced: true,
	}, api.NamespaceAll).Watch(&metav1.ListOptions{})

	if err != nil {
		return err
	}

	log.Printf("Watcher for %s created \n", types.IotDeviceType)

	defer watcher.Stop()

	for {
		e, ok := <-watcher.ResultChan()

		if !ok {
			return fmt.Errorf("%s watch ended due to a timeout", types.IotDeviceType)
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
			return fmt.Errorf("Error %s", types.IotDeviceType)
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
			err := kubernetes.DeletePod(w.restClient, pod)
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

func createTypeMeta(apiVersion string) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       types.IotPodKind,
		APIVersion: apiVersion,
	}
}
