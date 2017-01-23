package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
	"github.com/fest-research/iot-addon/pkg/apiserver/watch"
)

type IRawProxy interface {
	Put(*restful.Request) ([]byte, error)
	Get(*restful.Request) ([]byte, error)
	Post(*restful.Request) ([]byte, error)
	Patch(*restful.Request) ([]byte, error)
	Watch(*restful.Request) watch.Watcher
}

type RawProxy struct {
	serverAddress string
}

func NewRawProxy(serverAddress string) IRawProxy {
	return &RawProxy{serverAddress: serverAddress}
}

func (this RawProxy) Get(req *restful.Request) ([]byte, error) {
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)
	log.Printf("[Raw proxy] GET Request (%s)", requestPath)

	r, err := http.Get(requestPath)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[Raw proxy] GET Response (%s): %s", requestPath, string(body))
	return body, nil
}

func (this RawProxy) Put(req *restful.Request) ([]byte, error) {
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)

	defer req.Request.Body.Close()
	reqBody, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("[Raw proxy] PUT Request (%s): %s", requestPath, string(reqBody))

	r, err := http.NewRequest("PUT", requestPath, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[Raw proxy] PUT Response (%s): %s", requestPath, string(body))
	return body, nil
}

func (this RawProxy) Post(req *restful.Request) ([]byte, error) {
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)

	defer req.Request.Body.Close()
	reqBody, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("[Raw proxy] POST Request (%s): %s", requestPath, string(reqBody))

	r, err := http.Post(requestPath, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[Raw proxy] POST Response (%s): %s", requestPath, string(body))
	return body, nil
}

func (this RawProxy) Patch(req *restful.Request) ([]byte, error) {
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)

	defer req.Request.Body.Close()
	reqBody, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("[Raw proxy] PATCH Request (%s): %s", requestPath, string(reqBody))

	r, err := http.NewRequest("PATCH", requestPath, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[Raw proxy] PATCH Response (%s): %s", requestPath, string(body))
	return body, nil
}

func (this RawProxy) Watch(req *restful.Request) watch.Watcher {
	watcher := watch.NewRawWatcher()
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)

	go watcher.Watch(requestPath)
	return watcher
}

// Remove everything after '?' in url path (FOR TESTS ONLY!)
func (this RawProxy) removePathParams(url *url.URL) string {
	path := url.String()
	if strings.Contains(path, "?") {
		path = path[:strings.Index(path, "?")]
	}
	return path
}
