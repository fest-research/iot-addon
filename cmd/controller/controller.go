package main

import (
	"flag"

	"github.com/fest-research/iot-addon/pkg/controller/watch"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	"github.com/spf13/pflag"
)

var (
	apiserverArg  = pflag.String("apiserver", "", "apiserver adress in http://host:port format")
	kubeconfigArg = pflag.String("kubeconfig", "", "absolute path to the kubeconfig file")
)

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	flag.CommandLine.Parse(make([]string, 0))

	config := kubernetes.NewClientConfig(*apiserverArg, *kubeconfigArg)
	dynamicClient := kubernetes.NewDynamicClient(config)
	restClient := kubernetes.NewRESTClient(config)

	go watch.WatchIotDevices(dynamicClient)
	watch.WatchIotDaemonSet(dynamicClient, restClient)
}
