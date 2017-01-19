package watch

import (
	"fmt"
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"log"
)

var iotDaemonSetResource = v1.APIResource{
	Name:       "iotdaemonsets",
	Namespaced: true,
}

func WatchIotDaemonSet(client *dynamic.Client) {
	watcher, err := client.
		Resource(&iotDaemonSetResource, api.NamespaceAll).
		Watch(&api.ListOptions{})

	if err != nil {
		fmt.Println(err.Error())
	}

	defer watcher.Stop()

	for {
		e, ok := <-watcher.ResultChan()

		if !ok {
			panic(fmt.Sprintf("IotDaemonSet ended early?"))
		}

		iotDaemonSet, _ := e.Object.(*types.IotDaemonSet)

		if e.Type == watch.Added {
			log.Printf("Added %s\n", iotDaemonSet.Metadata.SelfLink)
			pod, _ := kubernetes.GetIotPod(*iotDaemonSet)
			log.Println(pod)

		} else if e.Type == watch.Modified {
			log.Printf("Modified %s\n", iotDaemonSet.Metadata.SelfLink)
		} else if e.Type == watch.Deleted {
			log.Printf("Deleted %s\n", iotDaemonSet.Metadata.SelfLink)
		} else if e.Type == watch.Error {
			log.Println("Error")
			break
		}
	}
}
