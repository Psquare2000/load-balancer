package health

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type HealthChecker struct {
	status map[string]bool
	mu     sync.RWMutex
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		status: make(map[string]bool),
	}
}

func (hc *HealthChecker) Start(urls []*url.URL, interval time.Duration) {
	go func() {
		for {
			for _, u := range urls {
				resp, err := http.Get(u.String())

				hc.mu.Lock()
				prev := hc.status[u.String()]
				var isHealthy bool

				if err != nil || (resp != nil && resp.StatusCode >= 400) {
					isHealthy = false
				} else {
					isHealthy = true
				}

				// Only log if health status has changed
				if isHealthy != prev {
					if isHealthy {
						log.Printf("[INFO] Backend %s has recovered\n", u.String())
					} else {
						log.Printf("[WARN] Backend %s is DOWN\n", u.String())
					}
				}

				hc.status[u.String()] = isHealthy
				hc.mu.Unlock()

				if resp != nil {
					resp.Body.Close()
				}
			}
			time.Sleep(interval)
		}
	}()
}

func (hc *HealthChecker) IsHealthy(u *url.URL) bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	return hc.status[u.String()]
}

func (hc *HealthChecker) GetAllStatuses() map[string]bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	snapshot := make(map[string]bool)
	for backend, status := range hc.status {
		snapshot[backend] = status
	}
	return snapshot
}
