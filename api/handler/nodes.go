package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"
)

func init() {
	// Create
	createHandler := &APIHandler{
		Path:           "/nodes",
		Parameters:     make([]*restful.Parameter, 0),
		HandlerFunc:    createNode,
		HTTPMethod:     "POST",
		ReturnedCode:   http.StatusOK,
		ReturnedMsg:    "OK",
		ReturnedObject: nil,
	}
	registerAPIHandler(createHandler)

	// Read
	getHandler := &APIHandler{
		Path:           "/nodes",
		Parameters:     make([]*restful.Parameter, 0),
		HandlerFunc:    getNode,
		HTTPMethod:     "GET",
		ReturnedCode:   http.StatusOK,
		ReturnedMsg:    "OK",
		ReturnedObject: nil,
	}
	registerAPIHandler(getHandler)

	// Update
	updateHandler := &APIHandler{
		Path:           "/nodes",
		Parameters:     make([]*restful.Parameter, 0),
		HandlerFunc:    updateNode,
		HTTPMethod:     "PUT",
		ReturnedCode:   http.StatusOK,
		ReturnedMsg:    "OK",
		ReturnedObject: nil,
	}
	registerAPIHandler(updateHandler)

	// Delete
	deleteHandler := &APIHandler{
		Path:           "/nodes",
		Parameters:     make([]*restful.Parameter, 0),
		HandlerFunc:    deleteNode,
		HTTPMethod:     "DELETE",
		ReturnedCode:   http.StatusOK,
		ReturnedMsg:    "OK",
		ReturnedObject: nil,
	}
	registerAPIHandler(deleteHandler)
}

func createNode(req *restful.Request, resp *restful.Response) {

}

func getNode(req *restful.Request, resp *restful.Response) {

}

func listNodes(req *restful.Request, resp *restful.Response) {

}

func updateNode(req *restful.Request, resp *restful.Response) {

}

func deleteNode(req *restful.Request, resp *restful.Response) {

}
