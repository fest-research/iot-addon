package watch

import (
	"log"
	"time"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

type IotDeviceWatcher struct {
	dynamicClient *dynamic.Client
	restClient    *rest.RESTClient
	clientset     *client.Clientset
	iotDomain     string
}

var iotDeviceResource = metav1.APIResource{
	Name:       types.IotDeviceType,
	Namespaced: true,
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

	var watcher watch.Interface = nil
	var err error = nil
	var resourceName string = types.TprIotDevice + "." + w.iotDomain
	ticker := time.NewTicker(time.Second * 4)
	defer ticker.Stop()

	for ok := true; ok; ok = watcher == nil {
		select {
		case <-ticker.C:
			watcher, err = w.dynamicClient.
				Resource(&iotDeviceResource, api.NamespaceAll).
				Watch(&metav1.ListOptions{})
			if err != nil {
				log.Println(err.Error())

				_, err = w.clientset.ExtensionsV1beta1().ThirdPartyResources().
					Get(resourceName, metav1.GetOptions{})

				if err != nil {
					tpr := &v1beta1.ThirdPartyResource{
						ObjectMeta: metav1.ObjectMeta{
							Name: resourceName,
						},
						Versions: []v1beta1.APIVersion{
							{Name: types.APIVersion},
						},
						Description: "A specification of a IoT Device",
					}

					_, err := w.clientset.ExtensionsV1beta1().ThirdPartyResources().Create(tpr)
					if err != nil {
						log.Println(err.Error())
					}
				}
			} else {
				ticker.Stop()
			}
			break
		}
	}

	defer watcher.Stop()
	log.Printf("Watcher for %s created \n", types.IotDeviceType)

	for {
		e := <-watcher.ResultChan()

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
		Body(&metav1.DeleteOptions{}).
		Do().
		Error()

}

func createTypeMeta(apiVersion string) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       types.IotPodKind,
		APIVersion: apiVersion,
	}
}
