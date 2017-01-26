package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/apiserver/controller"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"
	"github.com/fest-research/iot-addon/pkg/apiserver/watch"

	apimachinery "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/pkg/api"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

var iotDeviceResource = &apimachinery.APIResource{Name: v1.IotDeviceType, Namespaced: true}

type NodeService struct {
	proxy          proxy.IServerProxy
	nodeController controller.INodeController
}

// NewNodeService creates the API service for translating k8s Nodes into IotDevices.
func NewNodeService(proxy proxy.IServerProxy, controller controller.INodeController) NodeService {
	return NodeService{proxy: proxy, nodeController: controller}
}

// Register creates the api routes for the NodeService.
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
			Path("/nodes/{node}").
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
			Path("/nodes/{node}/status").
			To(this.updateStatus).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)
}

// TODO: refactor this method
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

	unstructuredIotDevice, err := this.proxy.Get(iotDeviceResource, namespace, node.Name)
	if err != nil {
		// Transform the node to an unstructured iot device
		unstructuredIotDevice, err = this.nodeController.ToUnstructured(node)
		if err != nil {
			handleInternalServerError(resp, err)
			return
		}

		// Create the iot device
		unstructuredIotDevice, err = this.proxy.Create(iotDeviceResource, namespace, unstructuredIotDevice)
		if err != nil {
			handleInternalServerError(resp, err)
			return
		}
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
	// TODO: refactor this later, set based on tenant
	namespace := "default"
	name := req.PathParameter("node")

	obj, err := this.proxy.Get(iotDeviceResource, namespace, name)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	response, err := this.nodeController.ToBytes(obj)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this NodeService) listNodes(req *restful.Request, resp *restful.Response) {
	// TODO: refactor this later, set based on tenant
	namespace := "default"

	obj, err := this.proxy.List(iotDeviceResource, namespace, &apiv1.ListOptions{})
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	iotDeviceList := obj.(*v1.IotDeviceList)
	nodeList := this.nodeController.ToNodeList(iotDeviceList)
	response, _ := json.Marshal(nodeList)

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this NodeService) updateStatus(req *restful.Request, resp *restful.Response) {
	// TODO: refactor this later, set based on tenant
	namespace := "default"
	name := req.PathParameter("node")

	// Read post request
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	// Update the IoTDevice
	unstructuredIotDevice, err := this.proxy.Patch(iotDeviceResource, namespace, name, api.MergePatchType, body)
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

func (this NodeService) watchNodes(req *restful.Request, resp *restful.Response) {
	// TODO: refactor this later, set based on tenant
	namespace := "default"

	watcher, err := this.proxy.Watch(iotDeviceResource, namespace, &apiv1.ListOptions{})
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	defer watcher.Stop()

	notifier := watch.NewNotifier()

	notifier.Register(this.nodeController)
	err = notifier.Start(watcher, resp)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}
}
