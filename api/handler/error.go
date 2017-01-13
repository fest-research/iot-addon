package handler

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
	"net/http"
)

func handleInternalServerError(response *restful.Response, err error) {
	log.Print(err)
	response.WriteError(http.StatusInternalServerError, err)
}