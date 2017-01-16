package proxy

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
)

type IServerProxy interface {
	Put(*restful.Request) ([]byte, error)
	Get(*restful.Request) ([]byte, error)
	Post(*restful.Request) ([]byte, error)
}

type ServerProxy struct {
	serverAddress string
}

func NewServerProxy(address string) ServerProxy {
	return ServerProxy{serverAddress: address}
}

func (this ServerProxy) Get(req *restful.Request) ([]byte, error) {
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)
	log.Printf("[Proxy] Making GET call to server (%s): %s", this.serverAddress, requestPath)

	r, err := http.Get(requestPath)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[Response filter] (%s) response: %s", requestPath, string(body))
	log.Printf("[Response filter] (%s) content type: %s", requestPath, r.Header.Get("Content-Type"))
	log.Printf("[Response filter] (%s) transfer encoding: %s", requestPath, r.Header.Get("Transfer-Encoding"))
	return body, nil
}

func (this ServerProxy) Put(req *restful.Request) ([]byte, error) {
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)
	log.Printf("Making PUT to server (%s): %s", this.serverAddress, requestPath)
	r, err := http.NewRequest("PUT", requestPath, req.Request.Body)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[Response filter] (%s) response: %s", requestPath, string(body))
	log.Printf("[Response filter] (%s) content type: %s", requestPath, r.Header.Get("Content-Type"))
	log.Printf("[Response filter] (%s) transfer encoding: %s", requestPath, r.Header.Get("Transfer-Encoding"))
	return body, nil
}

func (this ServerProxy) Post(req *restful.Request) ([]byte, error) {
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)
	log.Printf("Making POST to server (%s): %s", this.serverAddress, requestPath)

	r, err := http.Post(requestPath, "application/json", req.Request.Body)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[Response filter] (%s) response: %s", requestPath, string(body))
	log.Printf("[Response filter] (%s) content type: %s", requestPath, r.Header.Get("Content-Type"))
	log.Printf("[Response filter] (%s) transfer encoding: %s", requestPath, r.Header.Get("Transfer-Encoding"))
	return body, nil
}

// Remove everything after '?' in url path (FOR TESTS ONLY!)
func (this ServerProxy) removePathParams(url *url.URL) string {
	path := url.String()
	if strings.Contains(path, "?") {
		path = path[:strings.Index(path, "?")]
	}
	return path
}
