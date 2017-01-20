package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"

	kubeapi "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
)

type EventService struct {
	proxy proxy.IServerProxy
}

func NewEventService(proxy proxy.IServerProxy) EventService {
	return EventService{proxy: proxy}
}

func (this EventService) Register(ws *restful.WebService) {

	// Create event
	ws.Route(
		ws.Method("POST").
			Path("/namespaces/{namespace}/events").
			To(this.createEvent).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)

	// udate event
	ws.Route(
		ws.Method("PATCH").
			Path("/namespaces/{namespace}/events/{event}").
			To(this.updateEvent).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)
}


func (this EventService) createEvent(req *restful.Request, resp *restful.Response) {
	updateResponse, err := this.proxy.Post(req,kubeapi.APIResource{})
	if err != nil {
		handleInternalServerError(resp, err)
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(updateResponse)
}


func (this EventService) updateEvent(req *restful.Request, resp *restful.Response) {
	updateResponse, err := this.proxy.Patch(req,kubeapi.APIResource{})
	if err != nil {
		handleInternalServerError(resp, err)
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(updateResponse)
}