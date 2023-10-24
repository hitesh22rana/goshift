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
	servers := servers.Init()
	servers.Add("https://duckduckgo.com/", "https://www.github.com/")

	loadbalancer := loadbalancer.Init(servers)

	fmt.Printf("Loadbalancer listening on %s\n", url)

	http.HandleFunc("/", loadbalancer.ForwardRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
