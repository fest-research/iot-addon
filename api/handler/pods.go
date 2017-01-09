package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"
)

func init() {
	// Create
	createHandler := &APIHandler{
		Path:           "/{namespace}/pod",
		Parameters:     make([]*restful.Parameter, 0),
		HandlerFunc:    createPod,
		HTTPMethod:     "POST",
		ReturnedCode:   http.StatusOK,
		ReturnedMsg:    "OK",
		ReturnedObject: nil,
	}
	registerAPIHandler(createHandler)

	// Read
	getHandler := &APIHandler{
		Path:           "/{namespace}/pod",
		Parameters:     make([]*restful.Parameter, 0),
		HandlerFunc:    getPod,
		HTTPMethod:     "GET",
		ReturnedCode:   http.StatusOK,
		ReturnedMsg:    "OK",
		ReturnedObject: nil,
	}
	registerAPIHandler(getHandler)

	// Update
	updateHandler := &APIHandler{
		Path:           "/{namespace}/pod",
		Parameters:     make([]*restful.Parameter, 0),
		HandlerFunc:    updatePod,
		HTTPMethod:     "PUT",
		ReturnedCode:   http.StatusOK,
		ReturnedMsg:    "OK",
		ReturnedObject: nil,
	}
	registerAPIHandler(updateHandler)

	// Delete
	deleteHandler := &APIHandler{
		Path:           "/{namespace}/pod",
		Parameters:     make([]*restful.Parameter, 0),
		HandlerFunc:    deletePod,
		HTTPMethod:     "DELETE",
		ReturnedCode:   http.StatusOK,
		ReturnedMsg:    "OK",
		ReturnedObject: nil,
	}
	registerAPIHandler(deleteHandler)
}

func createPod(req *restful.Request, resp *restful.Response) {

}

func getPod(req *restful.Request, resp *restful.Response) {

}

func listPods(req *restful.Request, resp *restful.Response) {

}

func updatePod(req *restful.Request, resp *restful.Response) {

}

func deletePod(req *restful.Request, resp *restful.Response) {

}
