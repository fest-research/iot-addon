package watch

import (
	"log"

	types "github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/common"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"regexp"
)

var iotDeviceResource = metav1.APIResource{
	Name:       "iotdevices",
	Namespaced: true,
}

func WatchIotDevices(dynamicClient *dynamic.Client, restClient *rest.RESTClient) {
	watcher, err := dynamicClient.
		Resource(&iotDeviceResource, api.NamespaceAll).
		Watch(&api.ListOptions{})

	if err != nil {
		log.Println(err.Error())
	}

	defer watcher.Stop()

	for {
		e, ok := <-watcher.ResultChan()

		if !ok {
			panic("IotDevices ended early?")
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

func addDeviceHandler(restClient *rest.RESTClient, iotDevice types.IotDevice) {
	pods, _ := kubernetes.GetDevicePods(restClient, iotDevice)
	log.Printf("--Device pods %s %v\n", iotDevice.Metadata.Name, pods)
	log.Printf("--Device pods len %d\n", len(pods))

	daemonSets, _ := kubernetes.GetDeviceDaemonSets(restClient, iotDevice)
	log.Printf("--Device ds %s %v\n", iotDevice.Metadata.Name, daemonSets)
	log.Printf("--Device ds len %d\n", len(daemonSets))

	mapFromPods := createPodMapFromPods(pods)
	mapFromDs := createPodMapFromDaemonSets(daemonSets)

	deviceName := iotDevice.Metadata.Name

	for key, value := range mapFromDs {
		mapFromPodsValue, ok := mapFromPods[key]
		// Pod already exists
		if ok {
			// Pod different then daemon set. Update needed
			if !api.Semantic.DeepEqual(mapFromPodsValue, value) {
				//UPDATE
				log.Printf("Pods are not equal")
			} else {

				log.Printf("Pod %s for device %s already exist", mapFromPodsValue.Metadata.Name, deviceName)
			}

		} else { //Pod doesn't exist yet. Must be created

			err := createPod(restClient, value, deviceName)
			if err != nil {
				log.Printf(err.Error())
				continue
			}
			log.Printf("Created new pod %s ",
				mapFromPodsValue.Metadata.Name)
		}
	}

}

func createPodMapFromPods(pods []types.IotPod) map[string]types.IotPod {

	resultmap := map[string]types.IotPod{}
	// Extract name from pattern '{name}-{uuid}'
	var validID = regexp.MustCompile(`-[\w]{8}(-[\w]{4}){3}-[\w]{12}`)

	for _, item := range pods {

		name := validID.Split(item.Metadata.Name, 2)[0]

		pod := types.IotPod{
			TypeMeta: createTypeMeta(item.APIVersion),
			Metadata: createObjectMeta(name, item.Metadata.Namespace),
			Spec:     item.Spec,
		}
		resultmap[name] = pod
	}
	return resultmap
}

func createPodMapFromDaemonSets(deamonSets []types.IotDaemonSet) map[string]types.IotPod {
	resultmap := map[string]types.IotPod{}

	for _, item := range deamonSets {
		name := item.Metadata.Name

		pod := types.IotPod{
			TypeMeta: createTypeMeta(item.APIVersion),
			Metadata: createObjectMeta(name, item.Metadata.Namespace),
			Spec:     item.Spec.Template.Spec,
		}
		resultmap[name] = pod
	}
	return resultmap
}

func createPod(restClient *rest.RESTClient, pod types.IotPod, deviceName string) error {
	name := pod.Metadata.Name
	pod.Metadata.Name = name + "-" + string(common.NewUUID())
	pod.Metadata.Labels = map[string]string{
		types.CreatedBy:      types.IotDaemonSetType + "." + name,
		types.DeviceSelector: deviceName,
	}

	log.Printf("Create pod %v", pod)

	return restClient.Post().
		Namespace(pod.Metadata.Namespace).
		Resource(types.IotPodType).
		Body(&pod).
		Do().
		Error()

}

func createTypeMeta(apiVersion string) metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       "IotPod",
		APIVersion: apiVersion,
	}
}

func createObjectMeta(name string, namespace string) v1.ObjectMeta {
	return v1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
}
