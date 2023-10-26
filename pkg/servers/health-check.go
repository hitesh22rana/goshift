package servers

import (
	"fmt"
	"net/http"
	"os"
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

	if res.StatusCode != http.StatusOK {
		s.Health = false
		return
	}

	s.Health = true
}

func (s *ServersConfig) healthCheck() {
	for _, server := range s.Hosts {
		server.healthCheck()
	}
}

func StartHealthCheck(s *ServersConfig, interval int) {
	fmt.Println("INFO: Health check cron job started in the background")
	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Every(interval).Second().Do(s.healthCheck)
	if err != nil {
		fmt.Println("ERROR: Health check cron job failed to start")
		os.Exit(1)
	}

	scheduler.StartAsync()
}
