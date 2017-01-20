package watch

import (
	"fmt"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"log"
)

var iotDeviceResource = v1.APIResource{
	Name:       "iotdevices",
	Namespaced: true,
}

func WatchIotDevices(client *dynamic.Client) {
	watcher, err := client.
		Resource(&iotDeviceResource, api.NamespaceAll).
		Watch(&api.ListOptions{})

	if err != nil {
		fmt.Println(err.Error())
	}

	defer watcher.Stop()

	for {
		e, ok := <-watcher.ResultChan()

		if !ok {
			panic(fmt.Sprintf("IotDevices ended early?"))
		}

		iotDevice, _ := e.Object.(*types.IotDevice)

		if e.Type == watch.Added {
			log.Printf("Added %s\n", iotDevice.Metadata.SelfLink)
		} else if e.Type == watch.Modified {
			log.Printf("Modified %s\n", iotDevice.Metadata.SelfLink)
		} else if e.Type == watch.Deleted {
			log.Printf("Deleted %s\n", iotDevice.Metadata.SelfLink)
		} else if e.Type == watch.Error {
			log.Println("Error")
			break
		}
	}
}
