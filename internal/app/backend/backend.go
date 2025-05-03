package backend

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL    *url.URL
	alive  bool
	mux    sync.RWMutex
	Proxy  *httputil.ReverseProxy
	Weight uint
}

func NewBackend(url *url.URL, proxy *httputil.ReverseProxy, weight uint) *Backend {
	return &Backend{
		URL:    url,
		alive:  true,
		Proxy:  proxy,
		Weight: weight,
	}
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.alive
}
