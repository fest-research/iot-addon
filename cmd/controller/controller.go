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
	iotDomain     = pflag.String("domain", "fujitsu.com", "Domain name for IoT resources")
)

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	log.SetOutput(os.Stdout)
	log.Printf("IoT domain name %s", *iotDomain)

	config := kubernetes.NewClientConfig(*apiserverArg, *kubeconfigArg, *iotDomain)
	dynamicClient := kubernetes.NewDynamicClient(config)
	restClient := kubernetes.NewRESTClient(config)
	clientset := kubernetes.NewClientset(config)

	// Start IotDevices watch.
	go watch.NewIotDeviceWatcher(dynamicClient, restClient, clientset, *iotDomain).Watch()

	// Wait a second and start IotPods watch.
	time.Sleep(time.Second)
	go watch.NewIotPodWatcher(dynamicClient, restClient, clientset, *iotDomain).Watch()

	// Wait a second and start IotDaemonSet watch.
	time.Sleep(time.Second)
	watch.NewIotDaemonSetWatcher(dynamicClient, restClient, clientset, *iotDomain).Watch()

}
