package proxy

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/log"
	"net/http"
	"io/ioutil"
	"strings"
	"net/url"
)

type IServerProxy interface {
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
	log.Printf("Making call to server(%s): %s", this.serverAddress, requestPath)

	r, err := http.Get(requestPath)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("Server(%s) response: %s", this.serverAddress, string(body))
	return body, nil
}

func (this ServerProxy) Post(req *restful.Request) ([]byte, error) {
	requestPath := this.serverAddress + this.removePathParams(req.Request.URL)
	log.Printf("Making post to server(%s): %s", this.serverAddress, requestPath)

	r, err := http.Post(requestPath, "application/json", req.Request.Body)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("Server(%s) response: %s", this.serverAddress, string(body))
	return body, nil
}

// Remove everything after '?' in url path (FOR TESTS ONLY!)
func (this ServerProxy) removePathParams(url *url.URL) string {
	path := url.String()
	return path[:strings.Index(path, "?")]
}