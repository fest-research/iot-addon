package controller

import (
	"fmt"
	"reflect"

	"github.com/emicklei/go-restful/log"
	"k8s.io/apimachinery/pkg/watch"
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

func (this PodController) transform(in watch.Event) watch.Event {
	log.Printf("%v", in)
	log.Printf("%s", reflect.TypeOf(in.Object))

	return in
}

func NewPodController() *PodController {
	return &PodController{}
}
