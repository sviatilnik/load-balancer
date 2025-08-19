package algorithms

import (
	"load-balancer/internal/app/backend"
	"sync"
	"sync/atomic"
)

type WeightedRoundRobin struct {
	current   uint64
	callCount uint
	m         sync.Mutex
}

func (r *WeightedRoundRobin) Name() string {
	return "WeightedRoundRobin"
}

func (r *WeightedRoundRobin) GetNextBackend(backends []*backend.Backend) *backend.Backend {
	backendsCount := len(backends)

	if backendsCount == 0 {
		return nil
	}

	r.m.Lock()
	defer r.m.Unlock()

	next := r.current
	if r.callCount == 0 {
		next = atomic.AddUint64(&r.current, uint64(1)) % uint64(backendsCount)
		r.callCount = backends[int(next)].Weight
	}

	r.callCount--

	for i := 0; i < backendsCount; i++ {
		idx := (int(next) + i) % backendsCount

		if backends[idx].IsAlive() {
			return backends[idx]
		}
	}

	return nil
}
