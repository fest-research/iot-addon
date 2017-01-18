package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/IoT-apiserver/pkg/apiserver/proxy"
)

type NodeService struct {
	proxy proxy.IServerProxy
}

func NewNodeService(proxy proxy.IServerProxy) NodeService {
	return NodeService{proxy: proxy}
}

func (this NodeService) Register(ws *restful.WebService) {
	// Create node
	ws.Route(
		ws.Method("POST").
			Path("/nodes").
			To(this.createNode).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)

	// List nodes
	ws.Route(
		ws.Method("GET").
			Path("/nodes").
			To(this.listNodes).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)

	// Watch nodes
	ws.Route(
		ws.Method("GET").
			Path("/watch/nodes").
			To(this.watchNodes).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)
}

func (this NodeService) createNode(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.Post(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this NodeService) watchNodes(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.Get(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this NodeService) listNodes(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.Get(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}
