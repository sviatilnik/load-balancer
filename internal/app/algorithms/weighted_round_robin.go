package algorithms

import (
	"fmt"
	"load-balancer/internal/app/backend"
)

type WeightedRoundRobin struct {
	pool []*backend.Backend
}

func (r *WeightedRoundRobin) Name() string {
	return "WeightedRoundRobin"
}

func (r *WeightedRoundRobin) GetNextBackend(backends []*backend.Backend) *backend.Backend {
	backendsCount := len(backends)

	if backendsCount == 0 {
		return nil
	}

	if r.pool == nil || len(r.pool) == 0 {
		r.pool = make([]*backend.Backend, 0)
		for _, b := range backends {
			for i := 0; i < int(b.Weight); i++ {
				r.pool = append(r.pool, b)
			}
		}
	}
	fmt.Println(r.pool)

	var resultBackend *backend.Backend
	for i, b := range r.pool {
		if b.IsAlive() {
			resultBackend = b
			r.pool[i] = nil
			break
		}
	}

	filteredPool := make([]*backend.Backend, 0)
	for _, b := range r.pool {
		if b != nil {
			filteredPool = append(filteredPool, b)
		}
	}
	r.pool = filteredPool

	return resultBackend
}

func (r *WeightedRoundRobin) removeBackendFromPool(backends []*backend.Backend, index int) []*backend.Backend {
	return append(backends[:index], backends[index+1:]...)
}
