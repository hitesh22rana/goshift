package servers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
)

var client *http.Client = &http.Client{}

func (s *Server) healthCheck() {
	req, _ := http.NewRequest("HEAD", s.URL, nil)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Health check failed for", s.URL)
		s.Health = false
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		s.Health = false
		return
	}

	s.Health = true
}

func (s *ServersConfig) healthCheck() {
	for _, server := range s.List {
		server.healthCheck()
	}
}

func StartHealthCheck(s *ServersConfig) {
	fmt.Println("INFO: Health check cron job started in the background")
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(10).Second().Do(s.healthCheck)
	scheduler.StartAsync()
}
