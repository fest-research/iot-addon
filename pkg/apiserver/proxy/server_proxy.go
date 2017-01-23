package proxy

import (
	"github.com/emicklei/go-restful/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
)

type IServerProxy interface {
	Create(*metav1.APIResource, *unstructured.Unstructured, string) (*unstructured.Unstructured, error)
	Delete(*metav1.APIResource, string, *v1.DeleteOptions) error
	Patch(*metav1.APIResource, string, api.PatchType, []byte) (*unstructured.Unstructured, error)
	Update(*metav1.APIResource, *unstructured.Unstructured, string) (*unstructured.Unstructured, error)
	Get(*metav1.APIResource, string, string) (*unstructured.Unstructured, error)
	List(*metav1.APIResource, *api.ListOptions) (runtime.Object, error)
	Watch(*metav1.APIResource, *api.ListOptions) (watch.Interface, error)
}

type ServerProxy struct {
	// Third party resources client
	tprClient *dynamic.Client
}

func NewServerProxy(tprClient *dynamic.Client) IServerProxy {
	return &ServerProxy{tprClient: tprClient}
}

func (this ServerProxy) List(resource *metav1.APIResource, listOptions *api.ListOptions) (
	runtime.Object, error) {
	log.Printf("[Server proxy] LIST resource: %s, namespaced: %t", resource.Name, resource.Namespaced)

	return this.tprClient.Resource(resource, api.NamespaceAll).List(listOptions)
}

func (this ServerProxy) Get(resource *metav1.APIResource, namespace, name string) (
	*unstructured.Unstructured, error) {
	log.Printf("[Server proxy] GET resource: %s, namespaced: %t", resource.Name, resource.Namespaced)

	return this.tprClient.Resource(resource, namespace).Get(name)
}

func (this ServerProxy) Create(resource *metav1.APIResource, obj *unstructured.Unstructured, namespace string) (
	*unstructured.Unstructured, error) {
	log.Printf("[Server proxy] CREATE resource: %s, namespaced: %t", resource.Name, resource.Namespaced)

	return this.tprClient.Resource(resource, namespace).Create(obj)
}

func (this ServerProxy) Delete(resource *metav1.APIResource, name string,
	deleteOptions *v1.DeleteOptions) error {
	log.Printf("[Server proxy] DELETE resource: %s, namespaced: %t", resource.Name, resource.Namespaced)

	return this.tprClient.Resource(resource, api.NamespaceAll).Delete(name, deleteOptions)
}

func (this ServerProxy) Patch(resource *metav1.APIResource, name string, pt api.PatchType,
	body []byte) (*unstructured.Unstructured, error) {
	log.Printf("[Server proxy] PATCH resource: %s, namespaced: %t", resource.Name, resource.Namespaced)

	return this.tprClient.Resource(resource, api.NamespaceAll).Patch(name, pt, body)
}

func (this ServerProxy) Update(resource *metav1.APIResource, obj *unstructured.Unstructured,
	namespace string) (*unstructured.Unstructured, error) {
	log.Printf("[Server proxy] UPDATE resource: %v", obj)

	return this.tprClient.Resource(resource, namespace).Update(obj)
}

func (this ServerProxy) Watch(resource *metav1.APIResource, listOptions *api.ListOptions) (
	watch.Interface, error) {
	log.Printf("[Server proxy] WATCH resource: %s, namespaced: %t", resource.Name, resource.Namespaced)

	watcher, err := this.tprClient.
		Resource(resource, api.NamespaceAll).
		Watch(listOptions)

	return watcher, err
}
