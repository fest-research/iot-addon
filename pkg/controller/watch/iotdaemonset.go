package watch

import (
	"log"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
)

type IotDaemonSetWatcher struct {
	dynamicClient *dynamic.Client
	restClient    *rest.RESTClient
}

func NewIotDaemonSetWatcher(dynamicClient *dynamic.Client, restClient *rest.RESTClient) IotDaemonSetWatcher {
	return IotDaemonSetWatcher{dynamicClient: dynamicClient, restClient: restClient}
}

// WatchIotDaemonSet watches for IotDaemonSet events and handles them.
func (w IotDaemonSetWatcher) Watch() {
	watcher, err := w.dynamicClient.Resource(&v1.APIResource{
		Name:       types.IotDaemonSetType,
		Namespaced: true,
	}, api.NamespaceAll).Watch(&api.ListOptions{})

	if err != nil {
		log.Println(err.Error())
	}

	defer watcher.Stop()

	for {
		e, ok := <-watcher.ResultChan()

		if !ok {
			panic("IotDaemonSet ended early?")
		}

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
		if !kubernetes.IsPodCreated(w.restClient, ds, device) {
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
			kubernetes.UpdatePod(w.restClient, existingPod, ds.Spec.Template)
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
