package handler

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
	"github.com/fest-research/IoT-apiserver/api/proxy"

	"net/http"
	"time"
)

// nothing will ever be sent down this channel
var neverExitWatch <-chan time.Time = make(chan time.Time)

// timeoutFactory abstracts watch timeout logic for testing
type TimeoutFactory interface {
	TimeoutCh() (<-chan time.Time, func() bool)
}

// realTimeoutFactory implements timeoutFactory
type realTimeoutFactory struct {
	timeout time.Duration
}

// TimeoutChan returns a channel which will receive something when the watch times out,
// and a cleanup function to call when this happens.
func (w *realTimeoutFactory) TimeoutCh() (<-chan time.Time, func() bool) {
	if w.timeout == 0 {
		return neverExitWatch, func() bool { return false }
	}
	t := time.NewTimer(w.timeout)
	return t.C, t.Stop
}

type PodService struct {
	proxy proxy.IServerProxy
}

func NewPodService(proxy proxy.IServerProxy) PodService {
	return PodService{proxy: proxy}
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
	updateResponse, err := this.proxy.Put(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.Write(updateResponse)
}

func (this PodService) getPod(req *restful.Request, resp *restful.Response) {
	podResponse, err := this.proxy.Get(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(podResponse)
}

func (this PodService) listPods(req *restful.Request, resp *restful.Response) {
	response, err := this.proxy.Get(req)
	if err != nil {
		handleInternalServerError(resp, err)
	}

	resp.AddHeader("Content-Type", "application/json")
	resp.Write(response)
}

func (this PodService) watchPods(req *restful.Request, resp *restful.Response) {
	log.Print("Watch pods called")

	// ensure the connection times out
	//timeoutFactory := &realTimeoutFactory{5}
	//timeoutCh, cleanup := timeoutFactory.TimeoutCh()
	//defer cleanup()
	//
	//resp.Header().Set("Content-Type", "application/json;watch=stream")
	//resp.Header().Set("Transfer-Encoding", "chunked")
	//resp.WriteHeader(http.StatusOK)
	//
	//var unknown runtime.Unknown
	//internalEvent := &metav1.InternalEvent{}
	//buf := &bytes.Buffer{}
	//ch := s.Watching.ResultChan()
	//for {
	//	select {
	//	case <-cn.CloseNotify():
	//		return
	//	case <-timeoutCh:
	//		return
	//	case event, ok := <-ch:
	//		if !ok {
	//			// End of results.
	//			return
	//		}
	//
	//		obj := event.Object
	//		s.Fixup(obj)
	//		if err := s.EmbeddedEncoder.Encode(obj, buf); err != nil {
	//			// unexpected error
	//			utilruntime.HandleError(fmt.Errorf("unable to encode watch object: %v", err))
	//			return
	//		}
	//
	//	// ContentType is not required here because we are defaulting to the serializer
	//	// type
	//		unknown.Raw = buf.Bytes()
	//		event.Object = &unknown
	//
	//	// the internal event will be versioned by the encoder
	//		*internalEvent = metav1.InternalEvent(event)
	//		if err := e.Encode(internalEvent); err != nil {
	//			utilruntime.HandleError(fmt.Errorf("unable to encode watch object: %v (%#v)", err, e))
	//			// client disconnect.
	//			return
	//		}
	//		if len(ch) == 0 {
	//			flusher.Flush()
	//		}
	//
	//		buf.Reset()
	//	}
	//}
}
