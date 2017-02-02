package watch

import (
	"fmt"
	"log"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
)

type IotPodWatcher struct {
	dynamicClient *dynamic.Client
	restClient    *rest.RESTClient
	clientset     *client.Clientset
	iotDomain     string
}

func NewIotPodWatcher(dynamicClient *dynamic.Client, restClient *rest.RESTClient, clientset *client.Clientset,
	iotDomain string) IotPodWatcher {
	return IotPodWatcher{
		dynamicClient: dynamicClient,
		restClient:    restClient,
		clientset:     clientset,
		iotDomain:     iotDomain,
	}
}

func (w IotPodWatcher) Watch() {
	for {
		err := w.start()
		if err != nil {
			log.Printf("An error occured: %s", err.Error())
		}
	}
}

func (w IotPodWatcher) start() error {
	watcher, err := w.dynamicClient.Resource(&metav1.APIResource{
		Name:       types.IotPodType,
		Namespaced: true,
	}, api.NamespaceAll).Watch(&metav1.ListOptions{})

	if err != nil {
		return err
	}

	log.Printf("Watcher for %s created \n", types.IotPodType)

	defer watcher.Stop()

	for {
		e, ok := <-watcher.ResultChan()

		if !ok {
			return fmt.Errorf("%s watch ended due to a timeout", types.IotPodType)
		}

		iotPod, _ := e.Object.(*types.IotPod)

		if e.Type == watch.Deleted {
			log.Printf("Pod deleted %s\n", iotPod.Metadata.Name)

		} else if e.Type == watch.Error {
			return fmt.Errorf("Error %s", types.IotPodType)
		}
	}
}
