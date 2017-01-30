package handler

import (
	"github.com/fest-research/iot-addon/pkg/apiserver/controller"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"
)

type IServiceFactory interface {
	GetRegisteredServices() []IService
}

type ServiceFactory struct {
	proxy     *proxy.Proxy
	services  []IService
	iotDomain string
}

// NewServiceFactory creates a factory that registers all all supported services.
func NewServiceFactory(proxy *proxy.Proxy, iotDomain string) *ServiceFactory {
	factory := &ServiceFactory{proxy: proxy, services: make([]IService, 0)}
	factory.init()

	return factory
}

func (this *ServiceFactory) registerService(service IService) {
	this.services = append(this.services, service)
}

func (this *ServiceFactory) init() {
	// Version service
	this.registerService(NewVersionService(this.proxy.RawProxy))

	// Node service
	this.registerService(NewNodeService(this.proxy.ServerProxy, controller.NewNodeController(this.iotDomain)))

	// Pod service
	this.registerService(NewPodService(this.proxy.ServerProxy, controller.NewPodController(this.iotDomain)))

	// Event service
	this.registerService(NewEventService(this.proxy.RawProxy))

	// Kubernetes service
	this.registerService(NewKubeService(this.proxy.RawProxy))
}

// GetRegisteredServices returns the list of all API services that are currently registered.
func (this *ServiceFactory) GetRegisteredServices() []IService {
	return this.services
}
