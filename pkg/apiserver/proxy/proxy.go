package proxy

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"bytes"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
	"github.com/fest-research/iot-addon/pkg/apiserver/watch"
)

type IServerProxy interface {
	Put(*restful.Request) ([]byte, error)
	Get(*restful.Request) ([]byte, error)
	Post(*restful.Request) ([]byte, error)
	Watch(*restful.Request) watch.Watcher
}

type ServerProxy struct {
	serverAddress string
}

func NewServerProxy(address string) ServerProxy {
	return ServerProxy{serverAddress: address}
}

func (this ServerProxy) Get(req *restful.Request) ([]byte, error) {
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)
	log.Printf("[Proxy] GET Request (%s)", requestPath)

	r, err := http.Get(requestPath)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[Proxy] GET Response (%s): %s", requestPath, string(body))
	return body, nil
}

func (this ServerProxy) Put(req *restful.Request) ([]byte, error) {
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)

	defer req.Request.Body.Close()
	reqBody, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("[Proxy] PUT Request (%s): %s", requestPath, string(reqBody))

	r, err := http.NewRequest("PUT", requestPath, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[Proxy] PUT Response (%s): %s", requestPath, string(body))
	return body, nil
}

func (this ServerProxy) Post(req *restful.Request) ([]byte, error) {
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)

	defer req.Request.Body.Close()
	reqBody, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("[Proxy] POST Request (%s): %s", requestPath, string(reqBody))

	r, err := http.Post(requestPath, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[Proxy] POST Response (%s): %s", requestPath, string(body))
	return body, nil
}

func (this ServerProxy) Watch(req *restful.Request) watch.Watcher {
	watcher := watch.NewRawWatcher()
	// TODO map request path to third party resource watch path
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)

	go watcher.Watch(requestPath)
	return watcher
}

// Remove everything after '?' in url path (FOR TESTS ONLY!)
func (this ServerProxy) removePathParams(url *url.URL) string {
	path := url.String()
	if strings.Contains(path, "?") {
		path = path[:strings.Index(path, "?")]
	}
	return path
}
