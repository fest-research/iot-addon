package watch

import (
	"fmt"
	"net/http"
	"time"

	"reflect"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
	"github.com/fest-research/iot-addon/pkg/apiserver/controller"
)

var defaultTimeout = 10 * time.Minute

type Notifier struct {
	controllers []controller.Controller
	timeout     time.Duration
}

// Controllers are executed in registration order
func (this *Notifier) Register(controllers ...controller.Controller) {
	this.controllers = append(this.controllers, controllers...)
}

func (this *Notifier) SetTimeout(timeout time.Duration) {
	this.timeout = timeout
}

func (this *Notifier) Start(watcher Watcher, response *restful.Response) error {
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
	errorChan := watcher.ErrorChan()

	for {
		select {
		case <-cn.CloseNotify():
			return nil
		case <-timeoutCh:
			return nil
		case err := <-errorChan:
			return err
		case msg := <-resultChan:
			// Transform data if there are any transformers registered
			for _, controller := range this.controllers {
				transformed, err := controller.Transform(msg)
				if err != nil {
					return err
				}

				// We are expecting same type as we provided
				msg, ok = transformed.(string)
				if !ok {
					return fmt.Errorf("Transformation type mismatch. Provided: %s, got: %s",
						reflect.TypeOf(msg), reflect.TypeOf(transformed))
				}
			}

			log.Printf("[Notifier] Sending response to watch client: %s", msg)
			_, err := response.Write([]byte(msg))
			if err != nil {
				return err
			}

			if len(resultChan) == 0 {
				flusher.Flush()
			}
		}
	}
}

func NewNotifier() *Notifier {
	return &Notifier{controllers: make([]controller.Controller, 0), timeout: defaultTimeout}
}
