package handler

import restful "github.com/emicklei/go-restful"

type IService interface {
	Register(*restful.WebService)
}