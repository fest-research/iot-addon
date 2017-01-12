package handler

import (
	"net/http"

	"io/ioutil"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
)

const API_SERVER = "http://127.0.0.1:8080"

func init() {
	// Create
	createHandler := &APIHandler{
		Path:           "/pods",
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
		Path:           "/pods",
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
		Path:           "/pods",
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
		Path:           "/pods",
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
	log.Print(req)

	r, err := http.Get(API_SERVER + req.Request.URL.String())
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

func getPod(req *restful.Request, resp *restful.Response) {
	log.Print(req)

	r, err := http.Get(API_SERVER + req.Request.URL.String())
	if err != nil {
		log.Print(err)
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
	}
	log.Printf("API Server response: %s", string(body))
}

func listPods(req *restful.Request, resp *restful.Response) {
	log.Print(req)
}

func updatePod(req *restful.Request, resp *restful.Response) {
	log.Print(req)
}

func deletePod(req *restful.Request, resp *restful.Response) {
	log.Print(req)
}
