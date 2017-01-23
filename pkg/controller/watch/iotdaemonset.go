package watch

import (
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
	"log"
)

var iotDaemonSetResource = v1.APIResource{
	Name:       "iotdaemonsets",
	Namespaced: true,
}

func WatchIotDaemonSet(dynamicClient *dynamic.Client, restClient *rest.RESTClient) {
	watcher, err := dynamicClient.
		Resource(&iotDaemonSetResource, api.NamespaceAll).
		Watch(&api.ListOptions{})

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
			log.Printf("Added %s\n", ds.Metadata.SelfLink)
			kubernetes.CreateIotPods(*ds, dynamicClient, restClient)
			pods, _ := kubernetes.GetDaemonSetPods(restClient, *ds)
			log.Printf("Daemon set pods lenght is %d\n", len(pods))
		} else if e.Type == watch.Modified {
			log.Printf("Modified %s\n", ds.Metadata.SelfLink)
		} else if e.Type == watch.Deleted {
			log.Printf("Deleted %s\n", ds.Metadata.SelfLink)
		} else if e.Type == watch.Error {
			log.Println("Error")
			break
		}
	}
}
