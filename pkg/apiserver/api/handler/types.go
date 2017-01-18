package handler

import (
	"github.com/emicklei/go-restful"
)

type IService interface {
	Register(*restful.WebService)
}
