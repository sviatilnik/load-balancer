package algorithms

import (
	"fmt"
	"load-balancer/internal/app/backend"
	"sync/atomic"
)

type RoundRobin struct {
	current uint64
}

func (r *RoundRobin) Name() string {
	return "RoundRobin"
}

func (r *RoundRobin) GetNextBackend(backends []*backend.Backend) *backend.Backend {
	backendsCount := len(backends)

	next := atomic.AddUint64(&r.current, uint64(1)) % uint64(backendsCount)
	fmt.Println(next)
	fmt.Println(r.current)

	for i := 0; i < backendsCount; i++ {
		idx := (int(next) + i) % backendsCount

		fmt.Println(idx, backends[idx].IsAlive())
		if backends[idx].IsAlive() {
			return backends[idx]
		}
	}

	return nil
}
