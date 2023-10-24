package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var (
	servers []string = []string{
		"https://duckduckgo.com/",
		"https://www.github.com/",
	}
	port        int16 = 8000
	serverIndex int8  = 0
	mu          sync.Mutex
)

func getServer() *url.URL {
	mu.Lock()
	defer mu.Unlock()

	serverUrl, _ := url.Parse(servers[serverIndex])
	serverIndex = (serverIndex + 1) % int8(len(servers))
	return serverUrl
}

func forwardRequest(res http.ResponseWriter, req *http.Request) {
	url := getServer()
	fmt.Println("Forwarding request to", url)
	rProxy := httputil.NewSingleHostReverseProxy(url)
	rProxy.ServeHTTP(res, req)
}

func main() {
	fmt.Printf("Starting server on port %d\n", port)
	http.HandleFunc("/", forwardRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
