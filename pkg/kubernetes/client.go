package kubernetes

import (
	"log"

	"github.com/fest-research/iot-addon/pkg/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewDynamicClient(config *rest.Config) *dynamic.Client {
	client, err := dynamic.NewClient(config)

	if err != nil {
		panic(err.Error())
	}

	return client
}

func NewClientConfig(apiserver, kubeconfig string) *rest.Config {
	log.Printf("Creating client config using \"%s\" apiserver and \"%s\" kubeconfig",
		apiserver, kubeconfig)

	config, err := clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)

	if err != nil {
		panic(err.Error())
	}

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
	return config
}
