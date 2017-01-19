package handler

import (
	"net/http"

	"k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/iot-addon/pkg/apiserver/controller"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"
	"github.com/fest-research/iot-addon/pkg/apiserver/watch"
	"k8s.io/client-go/pkg/api"
)

type NodeService struct {
	proxy          proxy.IServerProxy
	nodeController *controller.NodeController
}

func NewNodeService(proxy proxy.IServerProxy, controller *controller.NodeController) NodeService {
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
	// TODO: add the correct resource type
	response, err := this.proxy.Post(req, v1.APIResource{})
	if err != nil {
		handleInternalServerError(resp, err)
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this NodeService) getNode(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.Get(req, v1.APIResource{})
	if err != nil {
		handleInternalServerError(resp, err)
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this NodeService) watchNodes(req *restful.Request, resp *restful.Response) {
	watcher, err := this.proxy.Watch(&v1.APIResource{}, &api.ListOptions{})
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	notifier := watch.NewNotifier()

	notifier.Register(this.nodeController)
	notifier.Start(watcher, resp)
}

func (this NodeService) listNodes(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.List(req, v1.APIResourceList{})
	if err != nil {
		handleInternalServerError(resp, err)
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this NodeService) updateStatus(req *restful.Request, resp *restful.Response) {
	updateResponse, err := this.proxy.Patch(req, v1.APIResource{})
	if err != nil {
		handleInternalServerError(resp, err)
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(updateResponse)
}
