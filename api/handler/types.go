package handler

import restful "github.com/emicklei/go-restful"

type APIHandler struct {
	Path           string
	Parameters     []*restful.Parameter
	HandlerFunc    restful.RouteFunction
	HTTPMethod     string
	ReturnedCode   int
	ReturnedMsg    string
	ReturnedObject interface{}
}
