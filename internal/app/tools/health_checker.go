package tools

import (
	"load-balancer/internal/app/backend"
	"log"
	"net"
	"time"
)

type HealthChecker struct {
	TimeOut  time.Duration
	Backends []*backend.Backend
}

func (h *HealthChecker) isBackendAlive(backend *backend.Backend) bool {
	conn, err := net.DialTimeout("tcp", backend.URL.Host, h.TimeOut)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

func (h *HealthChecker) Check() bool {
	result := true
	for _, back := range h.Backends {
		status := h.isBackendAlive(back)
		back.SetAlive(status)
		if status {
			log.Printf("backend %s is alive\n", back.URL.Host)
		} else {
			log.Printf("backend %s is not alive\n", back.URL.Host)
		}
		result = result && status
	}
	return result
}

func (h *HealthChecker) CheckWithPeriod(interval time.Duration) {
	t := time.NewTicker(interval)

	for {
		select {
		case <-t.C:
			h.Check()
		}
	}
}
