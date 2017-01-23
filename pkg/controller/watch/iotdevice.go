package watch

import (
	"fmt"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
	"log"
)

var iotDeviceResource = v1.APIResource{
	Name:       "iotdevices",
	Namespaced: true,
}

func WatchIotDevices(dynamicClient *dynamic.Client, restClient *rest.RESTClient) {
	watcher, err := dynamicClient.
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
			log.Printf("--Device added %s\n", iotDevice.Metadata.Name)
			addDeviceHandler(restClient, *iotDevice)
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

func addDeviceHandler (restClient *rest.RESTClient, iotDevice types.IotDevice) {
	pods, _ := kubernetes.GetDevicePods(restClient, iotDevice)
	log.Printf("--Device pods %s %v\n", iotDevice.Metadata.Name, pods)
	log.Printf("--Device pods len %d\n", len(pods))

	daemonSets, _ := kubernetes.GetDeviceDaemonSets(restClient, iotDevice)
	log.Printf("--Device ds %s %v\n", iotDevice.Metadata.Name, daemonSets)
	log.Printf("--Device ds len %d\n", len(daemonSets))

}