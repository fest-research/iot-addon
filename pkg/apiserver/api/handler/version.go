package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"
)

type VersionService struct {
	proxy proxy.IRawProxy
}

func NewVersionService(proxy proxy.IRawProxy) VersionService {
	return VersionService{proxy: proxy}
}

func (this VersionService) Register(ws *restful.WebService) {
	// Read
	ws.Route(
		ws.Method("GET").
			Path("/").
			To(this.getVersion).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)
}

func (this VersionService) getVersion(req *restful.Request, resp *restful.Response) {
	resp.Write([]byte(`{"Version": "v1"}`))
}
