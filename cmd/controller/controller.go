package main

import (
	"flag"

	"github.com/fest-research/iot-addon/pkg/controller/watch"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	"github.com/spf13/pflag"
	"time"
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

	// Start IotDevices watch.
	go watch.WatchIotDevices(dynamicClient, restClient)

	// Wait a second and start IotDaemonSet watch.
	time.Sleep(time.Second)
	watch.WatchIotDaemonSet(dynamicClient, restClient)
}
