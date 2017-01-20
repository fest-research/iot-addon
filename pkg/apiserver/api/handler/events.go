package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"
)

type EventService struct {
	proxy proxy.IRawProxy
}

func NewEventService(proxy proxy.IRawProxy) EventService {
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
	response, err := this.proxy.Post(req)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this EventService) updateEvent(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.Patch(req)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}
