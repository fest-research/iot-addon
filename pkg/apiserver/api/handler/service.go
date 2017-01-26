package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"
	"github.com/fest-research/iot-addon/pkg/apiserver/watch"
)

type KubeService struct {
	proxy proxy.IRawProxy
}

// NewKubeService creates the API service as a proxy to k8s Service resources.
func NewKubeService(proxy proxy.IRawProxy) KubeService {
	return KubeService{proxy: proxy}
}

// Register creates the API routes for the KubeService.
func (this KubeService) Register(ws *restful.WebService) {
	// List services
	ws.Route(
		ws.Method("GET").
			Path("/services").
			To(this.listServices).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)

	// Watch services
	ws.Route(
		ws.Method("GET").
			Path("/watch/services").
			To(this.watchServices).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)
}

func (this KubeService) listServices(req *restful.Request, resp *restful.Response) {

	response, err := this.proxy.Get(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this KubeService) watchServices(req *restful.Request, resp *restful.Response) {
	watcher := this.proxy.Watch(req)
	notifier := watch.NewRawNotifier()

	err := notifier.Start(watcher, resp)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}
}
