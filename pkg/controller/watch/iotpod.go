package watch

import (
	"log"
	"time"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

type IotPodWatcher struct {
	dynamicClient *dynamic.Client
	restClient    *rest.RESTClient
	clientset     *client.Clientset
	iotDomain     string
}

var iotPodResource = metav1.APIResource{
	Name:       types.IotPodType,
	Namespaced: true,
}

func NewIotPodWatcher(dynamicClient *dynamic.Client, restClient *rest.RESTClient, clientset *client.Clientset, iotDomain string) IotPodWatcher {
	return IotPodWatcher{dynamicClient: dynamicClient, restClient: restClient, clientset: clientset, iotDomain: iotDomain}
}

func (w IotPodWatcher) Watch() {

	var watcher watch.Interface = nil
	var err error = nil
	var resourceName string = types.TprIotPod+ "." + w.iotDomain
	ticker := time.NewTicker(time.Second * 4)
	defer ticker.Stop()

	for ok := true; ok; ok = watcher == nil {
		select {
		case <-ticker.C:
			watcher, err = w.dynamicClient.
				Resource(&iotPodResource, api.NamespaceAll).
				Watch(&metav1.ListOptions{})
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
						Description: "A specification of a IoT pod",
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
	log.Printf("Watcher for %s created \n", types.IotPodType)

	for {
		e := <-watcher.ResultChan()

		iotPod, _ := e.Object.(*types.IotPod)

		if e.Type == watch.Deleted {
			log.Printf("Pod deleted %s\n", iotPod.Metadata.Name)

		} else if e.Type == watch.Error {
			log.Println("Error")
			break
		}
	}
}
