package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"
)

func init() {
	// Read
	createHandler := &APIHandler{
		Path:           "/",
		Parameters:     make([]*restful.Parameter, 0),
		HandlerFunc:    getVersion,
		HTTPMethod:     "GET",
		ReturnedCode:   http.StatusOK,
		ReturnedMsg:    "OK",
		ReturnedObject: nil,
	}
	registerAPIHandler(createHandler)
}

func getVersion(req *restful.Request, resp *restful.Response) {
	resp.Write([]byte(`{"Version": "v1"}`))
}
