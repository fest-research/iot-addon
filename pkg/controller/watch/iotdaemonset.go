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

type IotDaemonSetWatcher struct {
	dynamicClient *dynamic.Client
	restClient    *rest.RESTClient
	clientset     *client.Clientset
	iotDomain     string
}

func NewIotDaemonSetWatcher(dynamicClient *dynamic.Client, restClient *rest.RESTClient, clientset *client.Clientset,
	iotDomain string) IotDaemonSetWatcher {
	return IotDaemonSetWatcher{
		dynamicClient: dynamicClient,
		restClient:    restClient,
		clientset:     clientset,
		iotDomain:     iotDomain,
	}
}

// WatchIotDaemonSet watches for IotDaemonSet events and handles them.
func (w IotDaemonSetWatcher) Watch() {
	var watcher watch.Interface = nil
	var err error = nil
	var resourceName string = types.TprIotDaemonSet + "." + w.iotDomain

	ticker := time.NewTicker(time.Second * 4)
	defer ticker.Stop()

	for ok := true; ok; ok = watcher == nil {
		select {
		case <-ticker.C:
			watcher, err = w.dynamicClient.Resource(&metav1.APIResource{
				Name:       types.IotDaemonSetType,
				Namespaced: true,
			}, api.NamespaceAll).Watch(&metav1.ListOptions{})
			if err != nil {
				log.Println(err.Error())
				_, err = w.clientset.ExtensionsV1beta1().ThirdPartyResources().Get(resourceName, metav1.GetOptions{})
				if err != nil {
					tpr := &v1beta1.ThirdPartyResource{
						ObjectMeta: metav1.ObjectMeta{
							Name: resourceName,
						},
						Versions: []v1beta1.APIVersion{
							{Name: types.APIVersion},
						},
						Description: "A specification of a IoT Daemon Set",
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
	log.Printf("Watcher for %s created \n", types.IotDaemonSetType)

	for {
		e := <-watcher.ResultChan()

		ds, _ := e.Object.(*types.IotDaemonSet)

		if e.Type == watch.Added {
			w.handleDaemonSetAddition(*ds)
		} else if e.Type == watch.Modified {
			w.handleDaemonSetModification(*ds)
		} else if e.Type == watch.Deleted {
			w.handleDaemonSetDeletion(*ds)
		} else if e.Type == watch.Error {
			log.Printf("Ending %s watch due to an error\n", types.IotDaemonSetType)
			break
		}
	}
}

// handleDaemonSetAddition handles new IotDaemonSet addition event. It creates IotPods for added IotDaemonSet if they
// don't exist yet.
func (w IotDaemonSetWatcher) handleDaemonSetAddition(ds types.IotDaemonSet) {
	log.Printf("Added new %s %s", types.IotDaemonSetKind, ds.Metadata.SelfLink)

	// Getting list of IotDaemonSet selected IotDevices.
	devices, err := kubernetes.GetDaemonSetDevices(ds, w.dynamicClient, w.restClient)
	if err != nil {
		log.Printf("Cannot get %s %s devices", types.IotDaemonSetKind, ds.Metadata.SelfLink)
	}

	// Creating IotPods on selected IotDevices if they don't exist yet.
	for _, device := range devices {

		unschedulable := kubernetes.GetUnschedulableLabelFromDevice(device)
		if !kubernetes.IsPodCreated(w.restClient, ds, device) && !unschedulable {
			kubernetes.CreateDaemonSetPod(ds, device, w.restClient)
		}
	}
}

// handleDaemonSetModification handles IotDaemonSet modification event. It reschedules IotPods for modified IotDaemonSet
// and updates their specs.
func (w IotDaemonSetWatcher) handleDaemonSetModification(ds types.IotDaemonSet) {
	log.Printf("Modified new %s %s", types.IotDaemonSetKind, ds.Metadata.SelfLink)

	// Making sure, that IotDaemonSet is deployed on currently selected IotDevices.
	// Getting all existing IotPods created by IotDaemonSet.
	existingPods, err := kubernetes.GetDaemonSetPods(w.restClient, ds)
	if err != nil {
		log.Printf("Cannot get %s %s pods", types.IotDaemonSetKind, ds.Metadata.SelfLink)
		return
	}

	// Getting list of IotDevices where IotDaemonSet should be deployed.
	destinedDevices, err := kubernetes.GetDaemonSetDevices(ds, w.dynamicClient, w.restClient)

	// Updating and deleting existing IotPods.
	for _, existingPod := range existingPods {
		if !kubernetes.IsPodCorrectlyScheduled(ds, existingPod) {
			kubernetes.DeletePod(w.restClient, existingPod)
		} else {
			err = kubernetes.UpdatePod(w.restClient, existingPod, ds.Spec.Template)
			if err != nil {
				log.Printf("Error. Can not update IotPod %s", existingPod.Metadata.Name)
			}
		}
	}

	// Add missing IotPods.
	for _, devicesMissingPod := range kubernetes.GetDevicesMissingPods(destinedDevices, existingPods) {
		kubernetes.CreateDaemonSetPod(ds, devicesMissingPod, w.restClient)
	}
}

// handleDaemonSetModification handles IotDaemonSet deletion event. It removes all IotPods created by deleted
// IotDaemonSet.
func (w IotDaemonSetWatcher) handleDaemonSetDeletion(ds types.IotDaemonSet) {
	log.Printf("Deleted %s %s", types.IotDaemonSetKind, ds.Metadata.SelfLink)

	// Deleting IotPods created by IotDaemonSet.
	kubernetes.DeleteDaemonSetPods(w.restClient, ds)
}
