package proxy

import (
	"github.com/emicklei/go-restful/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
)

type IServerProxy interface {
	Create(*metav1.APIResource, string, *unstructured.Unstructured) (*unstructured.Unstructured, error)
	Delete(*metav1.APIResource, string, string, *metav1.DeleteOptions) error
	Patch(*metav1.APIResource, string, string, types.PatchType, []byte) (*unstructured.Unstructured, error)
	Update(*metav1.APIResource, string, *unstructured.Unstructured) (*unstructured.Unstructured, error)
	Get(*metav1.APIResource, string, string) (*unstructured.Unstructured, error)
	List(*metav1.APIResource, string, *metav1.ListOptions) (runtime.Object, error)
	Watch(*metav1.APIResource, string, *metav1.ListOptions) (watch.Interface, error)
}

type ServerProxy struct {
	// Third party resources client
	tprClient *dynamic.Client
}

func NewServerProxy(tprClient *dynamic.Client) IServerProxy {
	return &ServerProxy{tprClient: tprClient}
}

func (this ServerProxy) Get(resource *metav1.APIResource, namespace, name string) (
	*unstructured.Unstructured, error) {
	log.Printf("[Server proxy] GET resource: %s, namespaced: %t", resource.Name, resource.Namespaced)

	return this.tprClient.Resource(resource, namespace).Get(name)
}

func (this ServerProxy) Create(resource *metav1.APIResource, namespace string, obj *unstructured.Unstructured) (
	*unstructured.Unstructured, error) {
	log.Printf("[Server proxy] CREATE resource: %s, namespaced: %t", resource.Name, resource.Namespaced)

	return this.tprClient.Resource(resource, namespace).Create(obj)
}

func (this ServerProxy) Delete(resource *metav1.APIResource, namespace, name string,
	deleteOptions *metav1.DeleteOptions) error {
	log.Printf("[Server proxy] DELETE resource: %s, namespaced: %t", resource.Name, resource.Namespaced)

	return this.tprClient.Resource(resource, namespace).Delete(name, deleteOptions)
}

func (this ServerProxy) Update(resource *metav1.APIResource, namespace string,
	obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	log.Printf("[Server proxy] UPDATE resource: %v", obj)

	return this.tprClient.Resource(resource, namespace).Update(obj)
}

func (this ServerProxy) Patch(resource *metav1.APIResource, namespace, name string,
	pt types.PatchType, body []byte) (*unstructured.Unstructured, error) {
	log.Printf("[Server proxy] PATCH resource: %s, namespaced: %t", resource.Name, resource.Namespaced)
	return this.tprClient.Resource(resource, namespace).Patch(name, pt, body)
}

func (this ServerProxy) List(resource *metav1.APIResource, namespace string, listOptions *metav1.ListOptions) (
	runtime.Object, error) {
	log.Printf("[Server proxy] LIST resource: %s, namespaced: %t", resource.Name, resource.Namespaced)
	return this.tprClient.Resource(resource, namespace).List(listOptions)
}

func (this ServerProxy) Watch(resource *metav1.APIResource, namespace string, listOptions *metav1.ListOptions) (
	watch.Interface, error) {
	log.Printf("[Server proxy] WATCH resource: %s, namespaced: %t", resource.Name, resource.Namespaced)

	watcher, err := this.tprClient.
		Resource(resource, namespace).
		Watch(listOptions)

	return watcher, err
}
