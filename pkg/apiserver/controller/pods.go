package controller

import (
	"fmt"
	"reflect"

	"github.com/fest-research/iot-addon/pkg/api/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"
	kubeapi "k8s.io/client-go/pkg/api/v1"
)

type Controller interface {
	// TODO: introduce into variable for more clear object transformation
	Transform(interface{}) (interface{}, error)
}

type PodController struct{}

func (this PodController) Transform(in interface{}) (interface{}, error) {
	switch in.(type) {
	case watch.Event:
		event := in.(watch.Event)
		return this.transformWatchEvent(event), nil
	case *v1.IotPodList:
		iotPodList := in.(*v1.IotPodList)
		return this.toPodList(iotPodList), nil
	case *v1.IotPod:
		iotPod := in.(*v1.IotPod)
		return this.toPod(iotPod), nil
	case *kubeapi.Pod:
		pod := in.(*kubeapi.Pod)
		return this.toUnstructured(pod)
	case *unstructured.Unstructured:
		unstructured := in.(*unstructured.Unstructured)
		return this.toBytes(unstructured)
	default:
		return nil, fmt.Errorf("Not supported type: %s", reflect.TypeOf(in))
	}
}

func (this PodController) transformWatchEvent(event watch.Event) watch.Event {
	iotPod := event.Object.(*v1.IotPod)
	event.Object = this.toPod(iotPod)
	return event
}

func (this PodController) toPodList(iotPodList *v1.IotPodList) kubeapi.PodList {
	podList := kubeapi.PodList{}

	podList.Kind = "PodList"
	podList.APIVersion = "v1"
	podList.Items = make([]kubeapi.Pod, 0)

	for _, iotPod := range iotPodList.Items {
		pod := this.toPod(&iotPod)
		podList.Items = append(podList.Items, *pod)
	}

	return podList
}

func (this PodController) toPod(iotPod *v1.IotPod) *kubeapi.Pod {
	pod := &kubeapi.Pod{}

	pod.Kind = "Pod"
	pod.APIVersion = "v1"
	pod.Spec = iotPod.Spec
	pod.ObjectMeta = iotPod.Metadata
	pod.Status = iotPod.Status

	pod.Spec.Containers[0].ImagePullPolicy = kubeapi.PullIfNotPresent
	pod.Spec.RestartPolicy = kubeapi.RestartPolicyAlways
	pod.Spec.DNSPolicy = kubeapi.DNSClusterFirst

	pod.Status.Phase = kubeapi.PodPending
	pod.Status.QOSClass = kubeapi.PodQOSBestEffort

	return pod
}

func (this PodController) toIotPod(pod *kubeapi.Pod) *v1.IotPod {
	iotPod := &v1.IotPod{}

	iotPod.Kind = "IotPod"
	iotPod.APIVersion = "fujitsu.com/v1"
	iotPod.Spec = pod.Spec
	iotPod.Metadata = pod.ObjectMeta
	iotPod.Status = pod.Status

	return iotPod
}

// Converts pod to unstructured iot pod
func (this PodController) toUnstructured(pod *kubeapi.Pod) (*unstructured.Unstructured, error) {
	result := &unstructured.Unstructured{}
	iotPod := this.toIotPod(pod)

	marshalledIotPod, err := json.Marshal(iotPod)
	if err != nil {
		return nil, err
	}

	err = result.UnmarshalJSON(marshalledIotPod)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Converts unstructured iot pod to pod json bytes array
func (this PodController) toBytes(unstructured *unstructured.Unstructured) ([]byte, error) {
	marshalledIotPod, err := unstructured.MarshalJSON()
	if err != nil {
		return nil, err
	}

	iotPod := &v1.IotPod{}
	err = json.Unmarshal(marshalledIotPod, iotPod)
	if err != nil {
		return nil, err
	}

	pod := this.toPod(iotPod)
	marshalledPod, err := json.Marshal(pod)
	if err != nil {
		return nil, err
	}

	return marshalledPod, nil
}

func NewPodController() *PodController {
	return &PodController{}
}
