package handler

import (
	"github.com/fest-research/iot-addon/pkg/apiserver/controller"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"
)

type IServiceFactory interface {
	GetRegisteredServices() []IService
}

type ServiceFactory struct {
	proxy    proxy.IServerProxy
	services []IService
}

func NewServiceFactory(proxy proxy.IServerProxy) *ServiceFactory {
	factory := &ServiceFactory{proxy: proxy, services: make([]IService, 0)}
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
	this.registerService(NewPodService(this.proxy, controller.NewPodController()))

	// Node service
	this.registerService(NewEventService(this.proxy))

	// Kubernetes service
	this.registerService(NewKubeService(this.proxy))
}

func (this *ServiceFactory) GetRegisteredServices() []IService {
	return this.services
}
