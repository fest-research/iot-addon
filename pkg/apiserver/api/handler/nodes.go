package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/iot-addon/pkg/apiserver/controller"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"
	"github.com/fest-research/iot-addon/pkg/apiserver/watch"
)

type NodeService struct {
	proxy          *proxy.Proxy
	nodeController *controller.NodeController
}

func NewNodeService(proxy *proxy.Proxy, controller *controller.NodeController) NodeService {
	return NodeService{proxy: proxy, nodeController: controller}
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

	// Get Node
	ws.Route(
		ws.Method("GET").
			Path("/nodes/{name}").
			To(this.getNode).
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

	// Update status
	ws.Route(
		ws.Method("PATCH").
			Path("/nodes/{name}/status").
			To(this.updateStatus).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)
}

func (this NodeService) createNode(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.RawProxy.Post(req)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this NodeService) getNode(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.RawProxy.Get(req)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this NodeService) watchNodes(req *restful.Request, resp *restful.Response) {
	watcher := this.proxy.RawProxy.Watch(req)
	notifier := watch.NewRawNotifier()

	err := notifier.Start(watcher, resp)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}
}

func (this NodeService) listNodes(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.RawProxy.Get(req)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this NodeService) updateStatus(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.RawProxy.Post(req)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}
