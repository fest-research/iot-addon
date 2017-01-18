package watch

import (
	"bufio"
	"bytes"
	"net/http"

	"github.com/emicklei/go-restful/log"
)

// Interface can be implemented by anything that knows how to watch and report changes.
// TODO: we might want to use some object type instead of chan string
type Watcher interface {
	Watch(string)

	ResultChan() chan string
	ErrorChan() chan error
}

type RawWatcher struct {
	result chan string
	err    chan error
}

// This is supposed to be called as go routine and synced using channels
func (this *RawWatcher) Watch(watchPath string) {
	log.Printf("[Watcher] Creating watch on %s", watchPath)
	resp, err := http.Get(watchPath)
	if err != nil {
		this.err <- err
		return
	}

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			this.err <- err
			return
		}

		log.Printf("[Watcher] Server response: %s", string(line))
		this.result <- bytes.NewBuffer(line).String()
	}
}

func (this *RawWatcher) ResultChan() chan string {
	return this.result
}

func (this *RawWatcher) ErrorChan() chan error {
	return this.err
}

func NewRawWatcher() Watcher {
	return &RawWatcher{result: make(chan string), err: make(chan error)}
}