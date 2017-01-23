package proxy

import "k8s.io/client-go/dynamic"

type Proxy struct {
	ServerProxy IServerProxy
	RawProxy    IRawProxy
}

func NewProxy(tprClient *dynamic.Client, serverAddress string) *Proxy {
	return &Proxy{ServerProxy: NewServerProxy(tprClient), RawProxy: NewRawProxy(serverAddress)}
}
