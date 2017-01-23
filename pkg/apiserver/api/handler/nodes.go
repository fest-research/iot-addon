package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/apiserver/controller"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"
	"github.com/fest-research/iot-addon/pkg/apiserver/watch"
	apimachinery "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

var iotDeviceResource = &apimachinery.APIResource{Name: v1.IotDeviceType, Namespaced: true}

type NodeService struct {
	proxy          *proxy.Proxy
	nodeController controller.INodeController
}

func NewNodeService(proxy *proxy.Proxy, controller controller.INodeController) NodeService {
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
	// TODO: refactor this later, set based on tenant
	namespace := "default"

	// Read post request
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	// Unmarshal request to a node object
	node := &apiv1.Node{}
	err = json.Unmarshal(body, node)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	// TODO: pass the namespace in Transform() when it's refactored
	node.ObjectMeta.Namespace = namespace

	// Transform the node to an unstructured iot device
	unstructuredIotDevice, err := this.nodeController.ToUnstructured(node)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	// Create the iot device
	unstructuredIotDevice, err = this.proxy.ServerProxy.Create(iotDeviceResource, unstructuredIotDevice, namespace)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	// Transform response back to unstructured pod
	response, err := this.nodeController.ToBytes(unstructuredIotDevice)
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
	response, err := this.proxy.RawProxy.Patch(req)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}
