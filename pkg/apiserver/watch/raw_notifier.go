package watch

import (
	"fmt"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
)

type RawNotifier struct {
	timeout time.Duration
}

func (this *RawNotifier) SetTimeout(timeout time.Duration) {
	this.timeout = timeout
}

func (this *RawNotifier) Start(watcher Watcher, response *restful.Response) error {
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
			log.Printf("[Raw Notifier] Sending response to watch client: %s", msg)
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

func NewRawNotifier() *RawNotifier {
	return &RawNotifier{timeout: defaultTimeout}
}
