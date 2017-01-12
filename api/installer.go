package api

import (
	restful "github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
	"github.com/fest-research/IoT-apiserver/api/handler"
)

// APIInstaller installs the APIs in the server
type APIInstaller struct {
	Root    string
	Version string
}

// NewWebService creates the core web service
func (installer *APIInstaller) NewWebService() *restful.WebService {
	ws := new(restful.WebService).Filter(webserviceLogging).Path(installer.Root).Consumes("*/*").Produces("application/json")
	ws.ApiVersion(installer.Version)
	return ws
}

// WebService Filter
func webserviceLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("[webservice-filter (logger)] \nRequest method: %s\nRequest path: %s\n", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
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
