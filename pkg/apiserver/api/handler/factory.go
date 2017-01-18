package handler

import (
	"github.com/fest-research/IoT-apiserver/pkg/apiserver/proxy"

	"k8s.io/client-go/kubernetes"
	"github.com/fest-research/IoT-apiserver/pkg/apiserver/controller"
)

type IServiceFactory interface {
	GetRegisteredServices() []IService
}

type ServiceFactory struct {
	clientSet *kubernetes.Clientset
	proxy     proxy.IServerProxy
	services  []IService
}

func NewServiceFactory(clientSet *kubernetes.Clientset, proxy proxy.IServerProxy) *ServiceFactory {
	factory := &ServiceFactory{clientSet: clientSet, proxy: proxy, services: make([]IService, 0)}
	factory.init()

	return factory
}

func (this *ServiceFactory) registerService(service IService) {
	this.services = append(this.services, service)
}

func (this *ServiceFactory) init() {
	// Version service
	this.registerService(NewVersionService(this.proxy))

	// Node service
	this.registerService(NewNodeService(this.proxy))

	// Pod service
	this.registerService(NewPodService(this.clientSet, this.proxy, controller.NewPodController()))

	// Kubernetes service
	this.registerService(NewKubeService(this.proxy))
}

func (this *ServiceFactory) GetRegisteredServices() []IService {
	return this.services
}
