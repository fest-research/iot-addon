package controller

import (
	"fmt"
	"reflect"

	"github.com/emicklei/go-restful/log"
)



type NodeController struct {}

func (this NodeController) Transform(in interface{}) (interface{}, error) {
	log.Print("NodeController - Transform()")

	switch in.(type) {
	case string:
		// TODO: do some transformation
		return in, nil
	default:
		return nil, fmt.Errorf("Not supported type: %s", reflect.TypeOf(in))
	}
}

func NewNodeController() *NodeController {
	return &NodeController{}
}