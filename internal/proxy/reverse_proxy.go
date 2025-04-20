package proxy

import (
	"go-load-balancer/internal/balancer"
	"net/http"
	"net/http/httputil"
)

type ProxyHandler struct {
	lb *balancer.RoundRobinBalancer
}

func NewProxyHandler(lb *balancer.RoundRobinBalancer) *ProxyHandler {
	return &ProxyHandler{lb: lb}
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := h.lb.NextBackend()
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}
