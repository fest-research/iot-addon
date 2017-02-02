package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/fest-research/iot-addon/pkg/controller/watch"
	"github.com/fest-research/iot-addon/pkg/kubernetes"
	"github.com/spf13/pflag"
)

var (
	apiserverArg  = pflag.String("apiserver", "", "apiserver adress in http://host:port format")
	kubeconfigArg = pflag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	iotDomain     = pflag.String("domain", "fujitsu.com", "custom domain name")
)

func main() {
	// Read command line arguments.
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	// Setup logger.
	log.SetOutput(os.Stdout)
	log.Printf("IoT domain name %s", *iotDomain)

	// Read cluster configuration.
	config := kubernetes.NewClientConfig(*apiserverArg, *kubeconfigArg, *iotDomain)

	// Create cluster clients.
	dynamicClient := kubernetes.NewDynamicClient(config)
	restClient := kubernetes.NewRESTClient(config)
	clientset := kubernetes.NewClientset(config)

	// Start watches.
	go watch.NewIotDeviceWatcher(dynamicClient, restClient, clientset, *iotDomain).Watch()
	go watch.NewIotPodWatcher(dynamicClient, restClient, clientset, *iotDomain).Watch()
	go watch.NewIotDaemonSetWatcher(dynamicClient, restClient, clientset, *iotDomain).Watch()

	// Avoid program exit.
	for {
		time.Sleep(time.Second)
	}
}
