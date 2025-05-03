package app

import (
	"load-balancer/internal/app/algorithms"
	"load-balancer/internal/app/backend"
	"log"
	"net/http"
)

type LoadBalancer struct {
	backends  []*backend.Backend
	algorithm algorithms.Algorithm
}

func NewLoadBalancer(algorithm algorithms.Algorithm, backends []*backend.Backend) *LoadBalancer {
	return &LoadBalancer{
		backends:  backends,
		algorithm: algorithm,
	}
}

func (l *LoadBalancer) AddBackend(back *backend.Backend) {
	l.backends = append(l.backends, back)
}

func (l *LoadBalancer) Backends() []*backend.Backend {
	return l.backends
}

func (l *LoadBalancer) SetBackends(backends []*backend.Backend) {
	l.backends = backends
}

func (l *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	back := l.algorithm.GetNextBackend(l.backends)
	if back == nil {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		log.Println("No available backends")
		return
	}

	back.Proxy.ServeHTTP(w, r)
}
