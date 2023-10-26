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
	getHealthyServer() (*servers.Server, error)
	ForwardRequest(res http.ResponseWriter, req *http.Request)
	Init(servers *servers.ServersConfig) LoadBalancerConfig
}

func (lb *LoadBalancerConfig) getHealthyServer() (*servers.Server, error) {
	for i := 0; i < len(lb.Hosts); i++ {
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
		fmt.Println("ERROR: No healthy servers found")
		res.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	fmt.Println("INFO: Forwarding request to", server.URL)
	server.ReverseProxy.ServeHTTP(res, req)
}

func Init(servers *servers.ServersConfig) LoadBalancerConfig {
	return LoadBalancerConfig{servers}
}
