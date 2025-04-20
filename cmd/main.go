package main

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"go-load-balancer/internal/balancer"
	"go-load-balancer/internal/health"
	"go-load-balancer/internal/proxy"
)

func main() {
	backendUrls := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	hc := health.NewHealthChecker()
	lb := balancer.NewRoundRobinBalancer(backendUrls, hc)
	handler := proxy.NewProxyHandler(lb)

	var parsedUrls []*url.URL
	for _, addr := range backendUrls {
		u, _ := url.Parse(addr)
		parsedUrls = append(parsedUrls, u)
	}
	hc.Start(parsedUrls, 5*time.Second)

	log.Println("[INFO] Load balancer running on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("[ERROR] Failed to start server:", err)
	}
}
