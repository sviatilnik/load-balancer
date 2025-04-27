package algorithms

import (
	"load-balancer/internal/app/backend"
)

type Algorithm interface {
	Name() string
	GetNextBackend(backends []*backend.Backend) *backend.Backend
}
