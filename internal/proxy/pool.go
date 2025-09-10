package proxy

import (
	"net/http"
	"sync"

	"api-gateway/config"
)

type backend struct {
	cfg   config.UpstreamConfig
	proxy http.Handler
	conns int
}

type Pool struct {
	mu       sync.RWMutex
	name     string
	strategy string
	backends []*backend
	rr       roundRobin
}

func NewPool(name, strategy string, backendsCfg []config.UpstreamConfig, m *Manager) (*Pool, error) {
	p := &Pool{name: name, strategy: strategy}
	for _, b := range backendsCfg {
		h, err := m.Get(b)
		if err != nil {
			return nil, err
		}
		p.backends = append(p.backends, &backend{cfg: b, proxy: h})
	}
	return p, nil
}

func (p *Pool) nextIdx() int {
	switch p.strategy {
	case "least_conn", "least-connections":
		conns := make([]int, len(p.backends))
		for i, b := range p.backends {
			conns[i] = b.conns
		}
		l := &leastConnections{}
		return l.Next(conns)
	default:
		return p.rr.Next(len(p.backends))
	}
}

func (p *Pool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mu.RLock()
	if len(p.backends) == 0 {
		p.mu.RUnlock()
		http.Error(w, "no backends", http.StatusBadGateway)
		return
	}
	idx := p.nextIdx()
	b := p.backends[idx]
	b.conns++
	p.mu.RUnlock()
	b.proxy.ServeHTTP(w, r)
	p.mu.Lock()
	b.conns--
	p.mu.Unlock()
}
