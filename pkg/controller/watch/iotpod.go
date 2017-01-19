package watch

import (
	"fmt"

	types "github.com/fest-research/IoT-apiserver/pkg/api/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
)

var iotPodResource = v1.APIResource{
	Name:       "iotpods",
	Namespaced: true,
}

func WatchIotPods(client *dynamic.Client) {
	watcher, err := client.
		Resource(&iotPodResource, api.NamespaceAll).
		Watch(&api.ListOptions{})

	if err != nil {
		fmt.Println(err.Error())
	}

	defer watcher.Stop()

	for {
		e, ok := <-watcher.ResultChan()

		if !ok {
			panic("IotPod ended early?")
		}

		iotPod, _ := e.Object.(*types.IotPod)

		if e.Type == watch.Added {
			fmt.Printf("Added %s\n", iotPod.Metadata.SelfLink)
		} else if e.Type == watch.Modified {
			fmt.Printf("Modified %s\n", iotPod.Metadata.SelfLink)
		} else if e.Type == watch.Deleted {
			fmt.Printf("Deleted %s\n", iotPod.Metadata.SelfLink)
		} else if e.Type == watch.Error {
			fmt.Println("Error")
			break
		}
	}
}
