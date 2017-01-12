package main

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/IoT-apiserver/api"
	"github.com/fest-research/IoT-apiserver/api/handler"
)

func main() {
	installer := api.APIInstaller{Root: "/api/v1", Version: "v1"}
	ws := installer.NewWebService()
	installer.Install(ws, handler.GetAPIHandlers())

	// DEBUG: check that the handlers are added properly
	fmt.Print(len(handler.GetAPIHandlers()))

	restful.Add(ws)
	routes := ws.Routes()
	fmt.Printf("%s", routes)
	http.ListenAndServe(":8083", nil)
}
