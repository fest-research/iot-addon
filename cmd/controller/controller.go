package main

import (
	"flag"
	"github.com/fest-research/iot-addon/pkg/controller/watch"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	"github.com/spf13/pflag"
	"k8s.io/client-go/dynamic"
	"log"
)

var (
	apiserverArg = pflag.String("apiserver", "",
		"Kubernetes API server host and port")
	kubeconfigArg = pflag.String("kubeconfig", "./kubeconfig.yaml",
		"absolute path to the kubeconfig file")
)

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	flag.CommandLine.Parse(make([]string, 0))

	if *apiserverArg != "" {
		log.Printf("Using apiserver location: %s\n", *apiserverArg)
	}

	if *kubeconfigArg != "" {
		log.Printf("Using kubeconfig file: %s\n", *kubeconfigArg)
	}

	config := kubernetes.GetConfig(*apiserverArg, *kubeconfigArg)
	client := kubernetes.NewDynamicClient(config)

	startWatching(client)
}

func startWatching(client *dynamic.Client) {
	go watch.WatchIotDevices(client)
	watch.WatchIotDaemonSet(client)
}
