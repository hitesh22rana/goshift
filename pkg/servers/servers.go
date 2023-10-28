package servers

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync"
)

var Mu sync.Mutex
var Wg sync.WaitGroup

type Server struct {
	URL          string
	ReverseProxy *httputil.ReverseProxy
	Health       bool
}

type ServersConfig struct {
	Hosts []*Server
	index int8
}

func (s *ServersConfig) shuffle() {
	Mu.Lock()
	defer Mu.Unlock()

	s.index = (s.index + 1) % int8(len(s.Hosts))
}

func (s *ServersConfig) Add(servers ...string) {
	for _, serverUrl := range servers {
		serverURL, err := url.Parse(serverUrl)
		if err != nil {
			fmt.Println("ERROR: URL parsing failed for", serverUrl)
		}

		server := Server{
			URL:          serverUrl,
			ReverseProxy: httputil.NewSingleHostReverseProxy(serverURL),
			Health:       false,
		}

		s.Hosts = append(s.Hosts, &server)
	}
}

func (s *ServersConfig) Current() *Server {
	server := s.Hosts[s.index]
	s.shuffle()
	return server
}

func Init() ServersConfig {
	return ServersConfig{}
}
