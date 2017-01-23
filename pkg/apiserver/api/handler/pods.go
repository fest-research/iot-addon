package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/fest-research/iot-addon/pkg/api/v1"
	"github.com/fest-research/iot-addon/pkg/apiserver/controller"
	"github.com/fest-research/iot-addon/pkg/apiserver/proxy"
	"github.com/fest-research/iot-addon/pkg/apiserver/watch"

	"fmt"

	apimachinery "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/pkg/api"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
)

var iotPodResource = &apimachinery.APIResource{Name: v1.IotPodType, Namespaced: true}

type PodService struct {
	proxy         proxy.IServerProxy
	podController controller.IPodController
}

func NewPodService(proxy proxy.IServerProxy, controller controller.IPodController) PodService {
	return PodService{proxy: proxy, podController: controller}
}

func (this PodService) Register(ws *restful.WebService) {
	// List pods
	ws.Route(
		ws.Method("GET").
			Path("/pods").
			To(this.listPods).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)

	// Watch pods
	ws.Route(
		ws.Method("GET").
			Path("/watch/pods").
			To(this.watchPods).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)

	// Get pod
	ws.Route(
		ws.Method("GET").
			Path("/namespaces/{namespace}/pods/{pod}").
			To(this.getPod).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)

	// Update pod status
	ws.Route(
		ws.Method("PUT").
			Path("/namespaces/{namespace}/pods/{pod}/status").
			To(this.updateStatus).
			Returns(http.StatusOK, "OK", nil).
			Writes(nil),
	)
}

func (this PodService) updateStatus(req *restful.Request, resp *restful.Response) {
	namespace := req.PathParameter("namespace")

	// Read update request
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	// Unmarshal request to pod object
	pod := &apiv1.Pod{}
	err = json.Unmarshal(body, pod)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	// Transform pod to unstructured iot pod
	unstructuredIotPod, err := this.podController.ToUnstructured(pod)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	// Update iot pod
	unstructuredIotPod, err = this.proxy.Update(iotPodResource, unstructuredIotPod, namespace)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	// Transform response back to unstructured pod
	response, err := this.podController.ToBytes(unstructuredIotPod)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this PodService) getPod(req *restful.Request, resp *restful.Response) {
	namespace := req.PathParameter("namespace")
	name := req.PathParameter("pod")

	obj, err := this.proxy.Get(iotPodResource, namespace, name)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	response, err := this.podController.ToBytes(obj)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this PodService) listPods(req *restful.Request, resp *restful.Response) {
	labelSelector, err := this.parseLabelSelector(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}

	fieldSelector, err := this.parseFieldSelector(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}
	obj, err := this.proxy.List(iotPodResource, &api.ListOptions{
		LabelSelector: *labelSelector,
		FieldSelector: *fieldSelector,
	})
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	iotPodList := obj.(*v1.IotPodList)
	podList := this.podController.ToPodList(iotPodList)
	response, _ := json.Marshal(podList)

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this PodService) watchPods(req *restful.Request, resp *restful.Response) {
	labelSelector, err := this.parseLabelSelector(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}

	fieldSelector, err := this.parseFieldSelector(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}

	watcher, err := this.proxy.Watch(iotPodResource, &api.ListOptions{
		LabelSelector: *labelSelector,
		FieldSelector: *fieldSelector,
	})
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}

	defer watcher.Stop()

	notifier := watch.NewNotifier()

	notifier.Register(this.podController)
	err = notifier.Start(watcher, resp)
	if err != nil {
		handleInternalServerError(resp, err)
		return
	}
}

func (this PodService) parseFieldSelector(req *restful.Request) (*fields.Selector, error) {
	selectorString := req.QueryParameter("fieldSelector")
	selector, err := fields.ParseSelector(selectorString)
	if err != nil {
		return nil, fmt.Errorf("[pod service] failed to parse field selector: %s", err)
	}
	return &selector, nil
}

func (this PodService) parseLabelSelector(req *restful.Request) (*labels.Selector, error) {
	selectorString := req.QueryParameter("labelSelector")
	selector, err := labels.Parse(selectorString)
	if err != nil {
		return nil, fmt.Errorf("[pod service] failed to parse label selector: %s", err)
	}
	return &selector, nil
}
