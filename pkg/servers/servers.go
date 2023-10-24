package servers

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

var Mu sync.Mutex

type Servers struct {
	list  []url.URL
	index int8
}

func handleError(err error, messagePrefix string) {
	if err != nil {
		fmt.Printf("%s: %s", messagePrefix, err)
		os.Exit(1)
	}
}

func (s *Servers) Add(servers ...string) {
	for _, server := range servers {
		serverUrl, err := url.Parse(server)
		handleError(err, "Error parsing server url")
		s.list = append(s.list, *serverUrl)
	}
}

func (s *Servers) GetCurrentServer() *url.URL {
	Mu.Lock()
	defer Mu.Unlock()

	serverUrl := s.list[s.index]
	s.index = (s.index + 1) % int8(len(s.list))
	return &serverUrl
}

func Init() Servers {
	return Servers{}
}
