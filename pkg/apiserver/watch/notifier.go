package watch

import (
	"fmt"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
	"github.com/fest-research/iot-addon/pkg/api/v1"
	ctrl "github.com/fest-research/iot-addon/pkg/apiserver/controller"

	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"
)

var defaultTimeout = 10 * time.Minute

type Notifier struct {
	controllers []ctrl.WatchEventController
	timeout     time.Duration
}

// Controllers are executed in registration order
func (this *Notifier) Register(controllers ...ctrl.WatchEventController) {
	this.controllers = append(this.controllers, controllers...)
}

func (this *Notifier) SetTimeout(timeout time.Duration) {
	this.timeout = timeout
}

func (this *Notifier) Start(watcher watch.Interface, response *restful.Response) error {
	log.Printf("[Notifier] Starting watch client notifier.")
	cn, ok := response.ResponseWriter.(http.CloseNotifier)
	if !ok {
		return fmt.Errorf("Unable to start watch - can't get http.CloseNotifier: %#v", response.ResponseWriter)
	}

	flusher, ok := response.ResponseWriter.(http.Flusher)
	if !ok {
		return fmt.Errorf("Unable to start watch - can't get http.Flusher: %#v", response.ResponseWriter)
	}

	// ensure the connection times out
	timeoutFactory := &realTimeoutFactory{timeout: this.timeout}
	timeoutCh, cleanup := timeoutFactory.TimeoutCh()
	defer cleanup()

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Transfer-Encoding", "chunked")
	response.WriteHeader(http.StatusOK)
	flusher.Flush()

	resultChan := watcher.ResultChan()

	for {
		select {
		case <-cn.CloseNotify():
			return nil
		case <-timeoutCh:
			return nil
		case event := <-resultChan:
			// Transform data if there are any controllers registered
			for _, controller := range this.controllers {
				event = controller.TransformWatchEvent(event)
			}

			// Our event has correct json annotations for watch event.
			iotEvent := this.toEvent(event)
			encodedEvent, err := json.Marshal(&iotEvent)
			if err != nil {
				return err
			}

			log.Printf("[Notifier] Sending response to watch client: %s", encodedEvent)
			_, err = response.Write(encodedEvent)
			if err != nil {
				return err
			}

			if len(resultChan) == 0 {
				flusher.Flush()
			}
		}
	}
}

func (this Notifier) toEvent(event watch.Event) v1.Event {
	return v1.Event{
		Type:   event.Type,
		Object: event.Object,
	}
}

func NewNotifier() *Notifier {
	return &Notifier{controllers: make([]ctrl.WatchEventController, 0), timeout: defaultTimeout}
}
