package algorithms

import (
	"load-balancer/internal/app/backend"
	"sync"
)

type WeightedRoundRobin struct {
	pool  []*backend.Backend
	mutex sync.RWMutex
}

func NewWeightedRoundRobin() *WeightedRoundRobin {
	return &WeightedRoundRobin{
		pool: make([]*backend.Backend, 0),
	}
}

func (r *WeightedRoundRobin) Name() string {
	return "WeightedRoundRobin"
}

func (r *WeightedRoundRobin) GetNextBackend(backends []*backend.Backend) *backend.Backend {
	/**
	TODO возможно не нужно всей этой работы с добавлением/удалением бекендов.
	Лучше при запуске создать сразу весь список бекендов и отправлять по ним так же как и в обычном RR алгоритме
	*/

	r.mutex.Lock()
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
	r.mutex.Unlock()

	return resultBackend
}

func (r *WeightedRoundRobin) removeBackendFromPool(backends []*backend.Backend, index int) []*backend.Backend {
	return append(backends[:index], backends[index+1:]...)
}
