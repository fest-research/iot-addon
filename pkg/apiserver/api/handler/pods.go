package handler

import (
	"net/http"
	"time"

	"k8s.io/client-go/kubernetes"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/IoT-apiserver/pkg/apiserver/controller"
	"github.com/fest-research/IoT-apiserver/pkg/apiserver/proxy"
	"github.com/fest-research/IoT-apiserver/pkg/apiserver/watch"
)

// nothing will ever be sent down this channel
var neverExitWatch <-chan time.Time = make(chan time.Time)

// timeoutFactory abstracts watch timeout logic for testing
type TimeoutFactory interface {
	TimeoutCh() (<-chan time.Time, func() bool)
}

// realTimeoutFactory implements timeoutFactory
type realTimeoutFactory struct {
	timeout time.Duration
}

// TimeoutChan returns a channel which will receive something when the watch times out,
// and a cleanup function to call when this happens.
func (w *realTimeoutFactory) TimeoutCh() (<-chan time.Time, func() bool) {
	if w.timeout == 0 {
		return neverExitWatch, func() bool { return false }
	}
	t := time.NewTimer(w.timeout)
	return t.C, t.Stop
}

type PodService struct {
	clientSet     *kubernetes.Clientset
	proxy         proxy.IServerProxy
	podController *controller.PodController
}

func NewPodService(clientSet *kubernetes.Clientset, proxy proxy.IServerProxy,
	controller *controller.PodController) PodService {
	return PodService{clientSet: clientSet, proxy: proxy, podController: controller}
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
	updateResponse, err := this.proxy.Put(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(updateResponse)
}

func (this PodService) getPod(req *restful.Request, resp *restful.Response) {
	podResponse, err := this.proxy.Get(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(podResponse)
}

func (this PodService) listPods(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.Get(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this PodService) watchPods(req *restful.Request, resp *restful.Response) {
	watcher := this.proxy.Watch(req)
	notifier := watch.NewNotifier()

	notifier.Register(this.podController)
	notifier.Start(watcher, resp)
}
