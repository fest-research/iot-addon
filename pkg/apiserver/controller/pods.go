package controller

import (
	"github.com/fest-research/iot-addon/pkg/api/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"
	kubeapi "k8s.io/client-go/pkg/api/v1"
)

type IPodController interface {
	// TransformWatchEvent implements WatchEventController.
	TransformWatchEvent(watch.Event) watch.Event

	ToPodList(*v1.IotPodList) *kubeapi.PodList
	ToPod(*v1.IotPod) *kubeapi.Pod
	ToIotPod(*kubeapi.Pod) *v1.IotPod
	ToUnstructured(*kubeapi.Pod) (*unstructured.Unstructured, error)
	ToBytes(*unstructured.Unstructured) ([]byte, error)
}

type podController struct{}

func (this podController) TransformWatchEvent(event watch.Event) watch.Event {
	iotPod := event.Object.(*v1.IotPod)
	event.Object = this.ToPod(iotPod)
	return event
}

func (this podController) ToPodList(iotPodList *v1.IotPodList) *kubeapi.PodList {
	podList := &kubeapi.PodList{}

	podList.Kind = "PodList"
	podList.APIVersion = "v1"
	podList.Items = make([]kubeapi.Pod, 0)

	for _, iotPod := range iotPodList.Items {
		pod := this.ToPod(&iotPod)
		podList.Items = append(podList.Items, *pod)
	}

	return podList
}

func (this podController) ToPod(iotPod *v1.IotPod) *kubeapi.Pod {
	pod := &kubeapi.Pod{}

	pod.Kind = "Pod"
	pod.APIVersion = "v1"
	pod.Spec = iotPod.Spec
	pod.ObjectMeta = iotPod.Metadata
	pod.Status = iotPod.Status

	for i := range pod.Spec.Containers {
		pod.Spec.Containers[i].ImagePullPolicy = kubeapi.PullAlways
	}
	pod.Spec.RestartPolicy = kubeapi.RestartPolicyAlways
	pod.Spec.DNSPolicy = kubeapi.DNSClusterFirst

	pod.Status.Phase = kubeapi.PodPending
	pod.Status.QOSClass = kubeapi.PodQOSBestEffort

	return pod
}

func (this podController) ToIotPod(pod *kubeapi.Pod) *v1.IotPod {
	iotPod := &v1.IotPod{}

	iotPod.Kind = "IotPod"
	iotPod.APIVersion = "fujitsu.com/v1"
	iotPod.Spec = pod.Spec
	iotPod.Metadata = pod.ObjectMeta
	iotPod.Status = pod.Status

	return iotPod
}

// Converts pod to unstructured iot pod
func (this podController) ToUnstructured(pod *kubeapi.Pod) (*unstructured.Unstructured, error) {
	result := &unstructured.Unstructured{}
	iotPod := this.ToIotPod(pod)

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
func (this podController) ToBytes(unstructured *unstructured.Unstructured) ([]byte, error) {
	marshalledIotPod, err := unstructured.MarshalJSON()
	if err != nil {
		return nil, err
	}

	iotPod := &v1.IotPod{}
	err = json.Unmarshal(marshalledIotPod, iotPod)
	if err != nil {
		return nil, err
	}

	pod := this.ToPod(iotPod)
	marshalledPod, err := json.Marshal(pod)
	if err != nil {
		return nil, err
	}

	return marshalledPod, nil
}

func NewPodController() IPodController {
	return &podController{}
}
