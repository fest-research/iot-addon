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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/pkg/api"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

var iotPodResource = &apimachinery.APIResource{Name: v1.IotPodType, Namespaced: true}

type PodService struct {
	proxy         proxy.IServerProxy
	podController *controller.PodController
}

func NewPodService(proxy proxy.IServerProxy, controller *controller.PodController) PodService {
	return PodService{proxy: proxy, podController: controller}
}

func (this PodService) Register(ws *restful.WebService) {
	// List pods
	ws.Route(
		ws.Method("GET").
			Path("/pods").
			To(this.listPods).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)

	// Watch pods
	ws.Route(
		ws.Method("GET").
			Path("/watch/pods").
			To(this.watchPods).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)

	// Get pod
	ws.Route(
		ws.Method("GET").
			Path("/namespaces/{namespace}/pods/{pod}").
			To(this.getPod).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)

	// Update pod status
	ws.Route(
		ws.Method("PUT").
			Path("/namespaces/{namespace}/pods/{pod}/status").
			To(this.updateStatus).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)
}

func (this PodService) updateStatus(req *restful.Request, resp *restful.Response) {
	namespace := req.PathParameter("namespace")

	// Read update request
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	// Unmarshal request to pod object
	pod := &apiv1.Pod{}
	err = json.Unmarshal(body, pod)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	// Transform pod to unstructured iot pod
	unstructuredIotPod, err := this.podController.Transform(pod)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	unstructured := unstructuredIotPod.(*unstructured.Unstructured)

	// Update iot pod
	unstructuredIotPod, err = this.proxy.Update(iotPodResource, unstructured, namespace)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	// Transform response back to unstructured pod
	response, err := this.podController.Transform(unstructuredIotPod)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	r := response.([]byte)

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(r)
}

func (this PodService) getPod(req *restful.Request, resp *restful.Response) {
	namespace := req.PathParameter("namespace")
	name := req.PathParameter("pod")

	obj, err := this.proxy.Get(iotPodResource, namespace, name)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	response, err := this.podController.Transform(obj)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	r := response.([]byte)

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(r)
}

func (this PodService) listPods(req *restful.Request, resp *restful.Response) {
	iotPodList, err := this.proxy.List(iotPodResource, &api.ListOptions{})
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	podListInterface, err := this.podController.Transform(iotPodList)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	podList := podListInterface.(apiv1.PodList)

	response, _ := json.Marshal(podList)

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this PodService) watchPods(req *restful.Request, resp *restful.Response) {
	watcher, err := this.proxy.Watch(iotPodResource, &api.ListOptions{})
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	defer watcher.Stop()

	notifier := watch.NewNotifier()

	notifier.Register(this.podController)
	err = notifier.Start(watcher, resp)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}
}
