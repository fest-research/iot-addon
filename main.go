package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/IoT-apiserver/api"
	"github.com/fest-research/IoT-apiserver/api/handler"
	"github.com/fest-research/IoT-apiserver/api/proxy"
	"github.com/spf13/pflag"
)

var (
	argApiserverHost = pflag.String("api-server", "", "Kubernetes api server address")
	argPort   = pflag.Int("port", 8083, "Port to listen on")
)

const (
	rootPath = "/api/v1"
	version  = "v1"
)

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	log.Printf("Using HTTP port: %d", *argPort)
	if *argApiserverHost == "" {
		log.Fatal("Parameter 'api-server' not defined. Please define kubernetes apiserver address.")
	}

	// Create api installer
	installer := api.APIInstaller{Root: rootPath, Version: version}

	// Create api proxy TODO: poll server and check if address is correct
	proxy := proxy.NewServerProxy(*argApiserverHost)

	// Create service factory
	serviceFactory := handler.NewServiceFactory(proxy)

	ws := installer.NewWebService()
	installer.Install(ws, serviceFactory.GetRegisteredServices())

	restful.Add(ws)
	http.ListenAndServe(fmt.Sprintf(":%d", *argPort), nil)
}
