package proxy

import (
	"go-load-balancer/internal/balancer"
	"log"
	"net/http"
	"net/http/httputil"
)

type ProxyHandler struct {
	lb         *balancer.RoundRobinBalancer
	maxRetries int
}

func NewProxyHandler(lb *balancer.RoundRobinBalancer) *ProxyHandler {
	return &ProxyHandler{
		lb:         lb,
		maxRetries: 3,
	}
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < h.maxRetries; i++ {
		target := h.lb.NextBackend()
		log.Printf("[INFO] Trying backend %s (attempt %d/%d)", target, i+1, h.maxRetries)

		proxy := httputil.NewSingleHostReverseProxy(target)

		// Track both proxy error and response code
		rw := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			log.Printf("[ERROR] Proxy to %s failed: %v", target, err)
		}

		proxy.ServeHTTP(rw, r)

		if rw.statusCode < 500 {
			return // success
		}

		// Retry if failed
		log.Printf("[WARN] Response status %d from %s — retrying", rw.statusCode, target)
	}

	http.Error(w, "All backends failed", http.StatusServiceUnavailable)
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}
