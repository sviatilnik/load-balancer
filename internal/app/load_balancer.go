package app

import (
	"load-balancer/internal/app/algorithms"
	"load-balancer/internal/app/backend"
	"net/http"
)

type LoadBalancer struct {
	backends  []*backend.Backend
	algorithm algorithms.Algorithm
}

func NewLoadBalancer(algorithm algorithms.Algorithm) *LoadBalancer {
	return &LoadBalancer{
		backends:  make([]*backend.Backend, 0),
		algorithm: algorithm,
	}
}

func (l *LoadBalancer) AddBackend(back *backend.Backend) {
	l.backends = append(l.backends, back)
}

func (l *LoadBalancer) Backends() []*backend.Backend {
	return l.backends
}

func (l *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	back := l.algorithm.GetNextBackend(l.backends)
	if back == nil {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	}

	back.Proxy.ServeHTTP(w, r)
}
