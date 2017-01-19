package controller

import (
	"fmt"
	"reflect"

	"github.com/emicklei/go-restful/log"
)

type Controller interface {
	Transform(interface{}) (interface{}, error)
}

type PodController struct{}

func (this PodController) Transform(in interface{}) (interface{}, error) {
	log.Print("PodController - Transform()")

	switch in.(type) {
	case string:
		// TODO: do some transformation
		return in, nil
	default:
		return nil, fmt.Errorf("Not supported type: %s", reflect.TypeOf(in))
	}
}

func NewPodController() *PodController {
	return &PodController{}
}
