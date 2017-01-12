package handler

import (
	"io/ioutil"
	"log"
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
	defer req.Request.Body.Close()
	log.Printf("Request header: %s", req.Request.Header)
	r, err := http.Post(API_SERVER+req.Request.URL.String(), "application/json", req.Request.Body)
	if err != nil {
		log.Print(err)
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}
	log.Printf("API Server response: %v", string(body))
}

func getNode(req *restful.Request, resp *restful.Response) {

}

func listNodes(req *restful.Request, resp *restful.Response) {

}

func updateNode(req *restful.Request, resp *restful.Response) {

}

func deleteNode(req *restful.Request, resp *restful.Response) {

}
