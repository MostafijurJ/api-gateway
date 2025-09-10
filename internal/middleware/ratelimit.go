package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"api-gateway/config"
)

type tokenBucket struct {
	mu     sync.Mutex
	tokens int
	last   time.Time
}

func (b *tokenBucket) allow(rate int, per time.Duration) bool {
	now := time.Now()
	elapsed := now.Sub(b.last)
	refill := int(float64(rate) * (float64(elapsed) / float64(per)))
	if refill > 0 {
		b.tokens = min(rate, b.tokens+refill)
		b.last = now
	}
	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

type limiterStore struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket
}

func newLimiterStore() *limiterStore { return &limiterStore{buckets: map[string]*tokenBucket{}} }

func (s *limiterStore) get(key string) *tokenBucket {
	s.mu.Lock()
	defer s.mu.Unlock()
	b, ok := s.buckets[key]
	if !ok {
		b = &tokenBucket{tokens: 0, last: time.Now()}
		s.buckets[key] = b
	}
	return b
}

func RateLimit(cfg config.RateConfig) Middleware {
	if !cfg.Enabled || cfg.Requests <= 0 || cfg.Per <= 0 {
		return func(next http.Handler) http.Handler { return next }
	}
	store := newLimiterStore()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			b := store.get(ip)
			if !b.allow(cfg.Requests, cfg.Per) {
				w.Header().Set("Retry-After", time.Duration(cfg.Per/time.Second).String())
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
