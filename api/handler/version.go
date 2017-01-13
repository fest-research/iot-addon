package handler

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/IoT-apiserver/api/proxy"
)

type VersionService struct {
	proxy proxy.IServerProxy
}

func NewVersionService(proxy proxy.IServerProxy) VersionService {
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
