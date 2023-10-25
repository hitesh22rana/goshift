package loadbalancer

import (
	"errors"
	"fmt"
	"net/http"

	servers "github.com/hitesh22rana/goshift/pkg/servers"
)

type LoadBalancerConfig struct {
	*servers.ServersConfig
}

type LoadBalancer interface {
	ForwardRequest(res http.ResponseWriter, req *http.Request)
}

func (lb *LoadBalancerConfig) getHealthyServer() (*servers.Server, error) {
	for i := 0; i < len(lb.List); i++ {
		server := lb.ServersConfig.Current()

		if server.Health {
			return server, nil
		}
	}

	return nil, errors.New("no healthy servers found")
}

func (lb *LoadBalancerConfig) ForwardRequest(res http.ResponseWriter, req *http.Request) {
	server, err := lb.getHealthyServer()

	if err != nil {
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	fmt.Println("INFO: Forwarding request to", server.URL)
	server.ReverseProxy.ServeHTTP(res, req)
}

func Init(servers *servers.ServersConfig) LoadBalancerConfig {
	return LoadBalancerConfig{servers}
}
