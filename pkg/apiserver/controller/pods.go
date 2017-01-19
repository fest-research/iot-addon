package controller

import (
	"fmt"
	"reflect"

	"github.com/emicklei/go-restful/log"
	"k8s.io/apimachinery/pkg/watch"
	"github.com/fest-research/iot-addon/pkg/api/v1"
)

type Controller interface {
	Transform(interface{}) (interface{}, error)
}

type PodController struct{}

func (this PodController) Transform(in interface{}) (interface{}, error) {
	log.Print("PodController - Transform()")

	switch in.(type) {
	case watch.Event:
		// TODO: do some transformation
		return in, nil
	default:
		return nil, fmt.Errorf("Not supported type: %s", reflect.TypeOf(in))
	}
}

func NewPodController() *PodController {
	return &PodController{}
}
