package handler

import (
	"github.com/fest-research/iot-addon/pkg/apiserver/controller"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"

	"k8s.io/client-go/dynamic"
)

type IServiceFactory interface {
	GetRegisteredServices() []IService
}

type ServiceFactory struct {
	kubeClient *dynamic.Client
	proxy      proxy.IServerProxy
	services   []IService
}

func NewServiceFactory(kubeClient *dynamic.Client, proxy proxy.IServerProxy) *ServiceFactory {
	factory := &ServiceFactory{kubeClient: kubeClient, proxy: proxy, services: make([]IService, 0)}
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
	this.registerService(NewPodService(this.kubeClient, this.proxy, controller.NewPodController()))

	// Kubernetes service
	this.registerService(NewKubeService(this.proxy))
}

func (this *ServiceFactory) GetRegisteredServices() []IService {
	return this.services
}
