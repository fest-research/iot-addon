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

func WatchIotDaemonSet(dynamicClient *dynamic.Client, restClient *rest.RESTClient) {
	watcher, err := dynamicClient.Resource(&v1.APIResource{
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
			handleDaemonSetAddition(dynamicClient, restClient, *ds)
		} else if e.Type == watch.Modified {
			handleDaemonSetModification(dynamicClient, restClient, *ds)
		} else if e.Type == watch.Deleted {
			handleDaemonSetDeletion(dynamicClient, restClient, *ds)
		} else if e.Type == watch.Error {
			log.Printf("Ending %s watch due to an error\n", types.IotDaemonSetType)
			break
		}
	}
}

func handleDaemonSetAddition(dynamicClient *dynamic.Client, restClient *rest.RESTClient, ds types.IotDaemonSet) {
	log.Printf("Added %s, creating pods...\n", ds.Metadata.SelfLink)
	kubernetes.CreateDaemonSetPods(ds, dynamicClient, restClient)
}

func handleDaemonSetModification(dynamicClient *dynamic.Client, restClient *rest.RESTClient, ds types.IotDaemonSet) {
	log.Printf("Modified %s\n", ds.Metadata.SelfLink)
	kubernetes.UpdateDaemonSetPods(restClient, ds)
}

func handleDaemonSetDeletion(dynamicClient *dynamic.Client, restClient *rest.RESTClient, ds types.IotDaemonSet) {
	log.Printf("Deleted %s\n", ds.Metadata.SelfLink)
	kubernetes.DeleteDaemonSetPods(restClient, ds)
}
