package watch

import (
	"fmt"
	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
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
			fmt.Printf("Added %s\n", iotDaemonSet.Metadata.SelfLink)
		} else if e.Type == watch.Modified {
			fmt.Printf("Modified %s\n", iotDaemonSet.Metadata.SelfLink)
		} else if e.Type == watch.Deleted {
			fmt.Printf("Deleted %s\n", iotDaemonSet.Metadata.SelfLink)
		} else if e.Type == watch.Error {
			fmt.Println("Error")
			break
		}
	}
}


