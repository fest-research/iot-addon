package controller

import (
	"github.com/fest-research/iot-addon/pkg/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type podController struct{
	iotDomain string
}

// TransformWatchEvent converts an ADD/UPDATE/DELETE event for an IotPod to
// an ADD/UPDATE/DELETE event for a k8s Pod
func (this podController) TransformWatchEvent(event watch.Event) watch.Event {
	iotPod := event.Object.(*v1.IotPod)
	event.Object = this.ToPod(iotPod)
	return event
}

// ToPodList converts a list of IotPods to a list of k8s Pods
func (this podController) ToPodList(iotPodList *v1.IotPodList) *kubeapi.PodList {
	podList := &kubeapi.PodList{}

	podList.TypeMeta = this.getTypeMeta(v1.PodListKind)
	podList.Items = make([]kubeapi.Pod, 0)

	for _, iotPod := range iotPodList.Items {
		pod := this.ToPod(&iotPod)
		podList.Items = append(podList.Items, *pod)
	}

	return podList
}

// ToPod converts an IotPod object to a k8s Pod object
func (this podController) ToPod(iotPod *v1.IotPod) *kubeapi.Pod {
	pod := &kubeapi.Pod{}

	pod.TypeMeta = this.getTypeMeta(v1.PodKind)

	pod.Spec = iotPod.Spec
	pod.ObjectMeta = iotPod.Metadata
	pod.Status = iotPod.Status

	pod = this.setRequiredFields(pod)

	return pod
}

// ToIotPod converts a k8s Pod object to an IotPod object
func (this podController) ToIotPod(pod *kubeapi.Pod) *v1.IotPod {
	iotPod := &v1.IotPod{}

	iotPod.TypeMeta = this.getIotTypeMeta()

	iotPod.Spec = pod.Spec
	iotPod.Metadata = pod.ObjectMeta
	iotPod.Status = pod.Status

	return iotPod
}

// ToUnstructured converts a k8s Pod object to an unstructured IotPod
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

// ToBytes converts an unstructured IotPod to a k8s Pod json bytes array
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

func (this podController) getIotTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: this.iotDomain,
		Kind:       v1.IotPodKind,
	}
}

func (this podController) getTypeMeta(kind v1.ResourceKind) metav1.TypeMeta {
	return metav1.TypeMeta{
		APIVersion: v1.APIVersion,
		Kind:       string(kind),
	}
}

func (this podController) setRequiredFields(pod *kubeapi.Pod) *kubeapi.Pod {
	for i := range pod.Spec.Containers {
		pod.Spec.Containers[i].ImagePullPolicy = kubeapi.PullAlways
		for j := range pod.Spec.Containers[i].Ports {
			pod.Spec.Containers[i].Ports[j].Protocol = kubeapi.ProtocolTCP
		}
	}

	pod.Spec.RestartPolicy = kubeapi.RestartPolicyAlways
	pod.Spec.DNSPolicy = kubeapi.DNSClusterFirst

	pod.Status.Phase = kubeapi.PodPending
	pod.Status.QOSClass = kubeapi.PodQOSBestEffort

	return pod
}

func NewPodController(iotDomain string) IPodController {
	return &podController{iotDomain: iotDomain}
}
