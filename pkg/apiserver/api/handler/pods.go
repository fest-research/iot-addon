package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/iot-addon/pkg/apiserver/controller"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"
	"github.com/fest-research/iot-addon/pkg/apiserver/watch"
	"k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/dynamic"
)

type PodService struct {
	kubeClient    *dynamic.Client
	proxy         proxy.IServerProxy
	podController *controller.PodController
}

func NewPodService(kubeClient *dynamic.Client, proxy proxy.IServerProxy,
	controller *controller.PodController) PodService {
	return PodService{kubeClient: kubeClient, proxy: proxy, podController: controller}
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
	updateResponse, err := this.proxy.Put(req, v1.APIResource{})
	if err != nil {
		handleInternalServerError(resp, err)
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(updateResponse)
}

func (this PodService) getPod(req *restful.Request, resp *restful.Response) {
	podResponse, err := this.proxy.Get(req, v1.APIResource{})
	if err != nil {
		handleInternalServerError(resp, err)
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(podResponse)
}

func (this PodService) listPods(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.List(req, v1.APIResourceList{})
	if err != nil {
		handleInternalServerError(resp, err)
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this PodService) watchPods(req *restful.Request, resp *restful.Response) {
	watcher := this.proxy.Watch(req, v1.APIResource{})
	notifier := watch.NewNotifier()

	notifier.Register(this.podController)
	notifier.Start(watcher, resp)
}
