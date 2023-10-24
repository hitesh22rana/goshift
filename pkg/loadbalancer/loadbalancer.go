package loadbalancer

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	servers "github.com/hitesh22rana/goshift/pkg/servers"
)

type loadbalancer struct {
	servers.Servers
}

func (lb *loadbalancer) ForwardRequest(res http.ResponseWriter, req *http.Request) {
	url := lb.Servers.GetCurrentServer()
	fmt.Println("Forwarding request to", url)
	rProxy := httputil.NewSingleHostReverseProxy(url)
	rProxy.ServeHTTP(res, req)
}

func Init(servers servers.Servers) loadbalancer {
	return loadbalancer{servers}
}
