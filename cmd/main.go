package main

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"encoding/json"
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
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		statuses := hc.GetAllStatuses()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statuses)
	})

	log.Println("[INFO] Load balancer running on :8080")
	// if err := http.ListenAndServe(":8080", handler); err != nil {
	// 	log.Fatal("[ERROR] Failed to start server:", err)
	// }
	http.Handle("/", handler)                    // default handler for all traffic
	log.Fatal(http.ListenAndServe(":8080", nil)) // nil = use default mux

}
