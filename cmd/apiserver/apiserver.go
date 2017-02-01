package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/apiserver/api"
	"github.com/fest-research/iot-addon/pkg/apiserver/api/handler"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"
	kube "github.com/fest-research/iot-addon/pkg/kubernetes"
	"github.com/spf13/pflag"
)

var (
	argApiserverHost = pflag.String("apiserver", "", "Kubernetes api server address")
	argPort          = pflag.Int("port", 8083, "Port to listen on")
	argKubeconfig    = pflag.String("kubeconfig", "", "Absolute path to the kubeconfig file")
	iotDomain        = pflag.String("domain", "fujitsu.com", "custom domain name")
)

const rootPath = "/api/" + v1.APIVersion

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	log.SetOutput(os.Stdout)
	log.Printf("Using HTTP port: %d", *argPort)
	if *argKubeconfig == "" && *argApiserverHost == "" {
		log.Println("Kubeconfig and apiserver arguments not provided. Falling back to inClusterConfig.")
	}

	// Get config object
	config := kube.NewClientConfig(*argApiserverHost, *argKubeconfig, *iotDomain)

	// Create a client for the kubernetes apis
	tprClient := kube.NewDynamicClient(config)

	// Create api installer
	installer := api.APIInstaller{Root: rootPath, Version: v1.APIVersion}

	// Create api proxy TODO: poll server and check if address is correct
	serverProxy := proxy.NewProxy(tprClient, config.Host)

	// Create service factory
	serviceFactory := handler.NewServiceFactory(serverProxy, *iotDomain)

	ws := installer.NewWebService()
	installer.Install(ws, serviceFactory.GetRegisteredServices())

	restful.Add(ws)
	http.ListenAndServe(fmt.Sprintf(":%d", *argPort), nil)
}
