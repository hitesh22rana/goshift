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
	List  []*Server
	index int8
}

type Servers interface {
	Add(servers ...string)
	Current() *Server
}

func (s *ServersConfig) shuffle() {
	Mu.Lock()
	defer Mu.Unlock()

	s.index = (s.index + 1) % int8(len(s.List))
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

		s.List = append(s.List, &server)
	}
}

func (s *ServersConfig) Current() *Server {
	server := s.List[s.index]
	s.shuffle()
	return server
}

func Init() ServersConfig {
	return ServersConfig{}
}
