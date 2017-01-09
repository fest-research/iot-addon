package api

import (
	restful "github.com/emicklei/go-restful"
	"github.com/taimir/small-kube/api/handler"
)

// APIInstaller installs the APIs in the server
type APIInstaller struct {
	Root    string
	Version string
}

// NewWebService creates the core web service
func (installer *APIInstaller) NewWebService() *restful.WebService {
	ws := new(restful.WebService).Path(installer.Root).Consumes("*/*").Produces("application/json")
	ws.ApiVersion(installer.Version)
	return ws
}

// Install installs the API handlers for all API resources
func (installer *APIInstaller) Install(ws *restful.WebService, apiHandlers []*handler.APIHandler) {
	for _, h := range apiHandlers {
		route := ws.Method(h.HTTPMethod).
			Path(h.Path).
			To(h.HandlerFunc).
			Returns(h.ReturnedCode, h.ReturnedMsg, h.ReturnedObject).
			Writes(h.ReturnedObject)
		for _, param := range h.Parameters {
			route.Param(param)
		}
		ws.Route(route)
	}
}
