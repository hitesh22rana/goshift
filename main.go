package main

import (
	"fmt"
	"log"
	"net/http"

	loadbalancer "github.com/hitesh22rana/goshift/pkg/loadbalancer"
	servers "github.com/hitesh22rana/goshift/pkg/servers"
)

var (
	port int    = 8000
	url  string = fmt.Sprintf("http://127.0.0.1:%d", port)
)

func main() {
	// Add your servers to the loadbalancer and start the healthcheck
	s := servers.Init()
	s.Add("https://api1.example.com", "https://api2.example.com")
	servers.StartHealthCheck(&s, 10)

	lb := loadbalancer.Init(&s)
	fmt.Printf("INFO: Loadbalancer listening on %s\n", url)

	// Forward the request to the loadbalancer
	http.HandleFunc("/", lb.ForwardRequest)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
