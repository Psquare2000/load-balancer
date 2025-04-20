package balancer

import (
	"go-load-balancer/internal/health"
	"log"
	"net/url"
	"sync/atomic"
)

type RoundRobinBalancer struct {
	backends []*url.URL
	counter  uint64
	health   *health.HealthChecker
}

func NewRoundRobinBalancer(urls []string, hc *health.HealthChecker) *RoundRobinBalancer {
	var backends []*url.URL
	for _, addr := range urls {
		parsed, err := url.Parse(addr)
		if err != nil {
			log.Fatalf("[ERROR] Invalid backend URL: %s", addr)
		}
		backends = append(backends, parsed)
	}
	return &RoundRobinBalancer{
		backends: backends,
		health:   hc,
	}
}

func (rr *RoundRobinBalancer) NextBackend() *url.URL {
	for i := 0; i < len(rr.backends); i++ {
		idx := atomic.AddUint64(&rr.counter, 1)
		candidate := rr.backends[idx%uint64(len(rr.backends))]
		if rr.health == nil || rr.health.IsHealthy(candidate) {
			return candidate
		}
	}
	return rr.backends[0] // fallback if all are unhealthy
}
