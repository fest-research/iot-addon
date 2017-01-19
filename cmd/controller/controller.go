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
	client := kubernetes.NewDynamicClient(config)

	go watch.WatchIotDevices(client)
	watch.WatchIotDaemonSet(client)
}
