package main

import (
	"flag"
	"log"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/fest-research/IoT-apiserver/pkg/api/v1"
	"github.com/fest-research/IoT-apiserver/pkg/controller/watch"
	"github.com/spf13/pflag"
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

	config, err := clientcmd.BuildConfigFromFlags(*apiserverArg, *kubeconfigArg)
	if err != nil {
		panic(err.Error())
	}

	configureClient(config)

	client, err := dynamic.NewClient(config)

	go watch.WatchIotDevices(client)

	watch.WatchIotDaemonSet(client)

	for {
		// Endless loop
	}
}

func configureClient(config *rest.Config) {
	groupVersion := schema.GroupVersion{
		Group:   "fujitsu.com",
		Version: "v1",
	}

	config.GroupVersion = &groupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: api.Codecs}

	schemeBuilder := runtime.NewSchemeBuilder(
		func(scheme *runtime.Scheme) error {
			scheme.AddKnownTypes(
				groupVersion,
				&api.ListOptions{},
				&api.DeleteOptions{},
				&v1.IotDevice{},
				&v1.IotDeviceList{},
				&v1.IotDaemonSet{},
				&v1.IotDaemonSetList{},
				&v1.IotPod{},
				&v1.IotPodList{},
			)
			return nil
		})

	schemeBuilder.AddToScheme(api.Scheme)
}
