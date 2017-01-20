package handler

import (
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
	"github.com/emicklei/go-restful/log"
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
	obj, err := this.proxy.Patch(&apimachinery.APIResource{}, "", api.JSONPatchType, []byte{})
	if err != nil {
		handleInternalServerError(resp, err)
	}

	response, err := obj.MarshalJSON()
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this PodService) getPod(req *restful.Request, resp *restful.Response) {
	obj, err := this.proxy.Get(&apimachinery.APIResource{}, "")
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	response, err := obj.MarshalJSON()
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
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
	log.Printf("WATCH PODS")
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
