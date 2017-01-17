package api

import (
	restful "github.com/emicklei/go-restful"
	//"github.com/emicklei/go-restful/log"
	"github.com/fest-research/IoT-apiserver/api/handler"
)

// APIInstaller installs the APIs in the server
type APIInstaller struct {
	Root    string
	Version string
}

// NewWebService creates the core web service
func (installer *APIInstaller) NewWebService() *restful.WebService {
	ws := new(restful.WebService).Filter(logPath).Path(installer.Root).Consumes("*/*").Produces("application/json")
	ws.ApiVersion(installer.Version)
	return ws
}

// Install installs the API handlers for all API resources
func (installer *APIInstaller) Install(ws *restful.WebService, services []handler.IService) {
	for _, s := range services {
		s.Register(ws)
	}
}

func logPath(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
	//log.Printf("[Request filter] method: %s ", req.Request.Method)
	//log.Printf("[Request filter] path: %s", req.Request.URL.String())
	//log.Printf("[Request filter] content type: %s", req.HeaderParameter("Content-Type"))
	//log.Printf("[Request filter] transfer encoding: %s", req.HeaderParameter("Transfer-Encoding"))
	chain.ProcessFilter(req, res)
}
