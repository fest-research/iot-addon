package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
)

func handleInternalServerError(response *restful.Response, err error) {
	log.Print(err)
	response.WriteError(http.StatusInternalServerError, err)
}
