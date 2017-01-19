package controller

import (
	"fmt"
	"reflect"

	"github.com/emicklei/go-restful/log"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/fest-research/iot-addon/pkg/api/v1"
	kubeapi "k8s.io/client-go/pkg/api/v1"
)

type Controller interface {
	Transform(interface{}) (interface{}, error)
}

type PodController struct{}

func (this PodController) Transform(in interface{}) (interface{}, error) {
	log.Print("PodController - Transform()")

	switch in.(type) {
	case watch.Event:
		event := in.(watch.Event)
		return this.transform(event), nil
	default:
		return nil, fmt.Errorf("Not supported type: %s", reflect.TypeOf(in))
	}
}

func (this PodController) transform(event watch.Event) watch.Event {
	log.Printf("%v", event)
	log.Printf("%s", reflect.TypeOf(event.Object))

	iotPod := event.Object.(*v1.IotPod)
	pod := kubeapi.Pod{}

	pod.Kind = "Pod"
	pod.APIVersion = "v1"
	pod.Spec = iotPod.Spec
	pod.ObjectMeta = iotPod.Metadata
	pod.Status = iotPod.Status

	event.Object = &pod
	return event
}

func NewPodController() *PodController {
	return &PodController{}
}
