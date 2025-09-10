package proxy

import (
	"net/http"
	"sync"

	"api-gateway/config"
)

type Manager struct {
	mu         sync.RWMutex
	urlToProxy map[string]http.Handler
}

func NewManager() *Manager { return &Manager{urlToProxy: make(map[string]http.Handler)} }

func (m *Manager) Get(cfg config.UpstreamConfig) (http.Handler, error) {
	m.mu.RLock()
	h, ok := m.urlToProxy[cfg.URL]
	m.mu.RUnlock()
	if ok {
		return h, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if h, ok := m.urlToProxy[cfg.URL]; ok {
		return h, nil
	}
	rp, err := NewReverseProxy(cfg)
	if err != nil {
		return nil, err
	}
	m.urlToProxy[cfg.URL] = rp
	return rp, nil
}

func (m *Manager) NewPool(name, strategy string, backends []config.UpstreamConfig) (*Pool, error) {
	return NewPool(name, strategy, backends, m)
}
