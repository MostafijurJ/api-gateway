package server

import (
	"net/http"
	"strings"
	"time"

	"api-gateway/config"
	"api-gateway/internal/middleware"
	"api-gateway/internal/observability"
	"api-gateway/internal/proxy"
	"api-gateway/internal/router"
)

func New(cfg config.Config) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.Handle("/metrics", observability.MetricsHandler())

	r := router.New(cfg.Routes)
	pm := proxy.NewManager()
	pools := map[string]*proxy.Pool{}
	for _, pc := range cfg.Pools {
		p, err := pm.NewPool(pc.Name, pc.Strategy, pc.Backends)
		if err != nil {
			observability.Logger().Error().Err(err).Str("pool", pc.Name).Msg("failed to init pool")
			continue
		}
		pools[pc.Name] = p
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		match := r.Match(req)
		up := cfg.Upstream
		strip := false
		if match != nil {
			if match.Pool != "" {
				if p, ok := pools[match.Pool]; ok {
					if match.StripPrefix && strings.HasPrefix(req.URL.Path, match.PathPrefix) {
						req.URL.Path = strings.TrimPrefix(req.URL.Path, match.PathPrefix)
						if req.URL.Path == "" {
							req.URL.Path = "/"
						}
					}
					p.ServeHTTP(w, req)
					return
				}
			}
			if match.Upstream.URL != "" {
				up = match.Upstream
			}
			strip = match.StripPrefix
		}
		h, err := pm.Get(up)
		if err != nil {
			observability.Logger().Error().Err(err).Msg("failed to get reverse proxy")
			http.Error(w, "bad gateway", http.StatusBadGateway)
			return
		}
		if strip && match != nil && strings.HasPrefix(req.URL.Path, match.PathPrefix) {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, match.PathPrefix)
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
		}
		h.ServeHTTP(w, req)
	})

	skipAuth := func(r *http.Request) bool {
		p := r.URL.Path
		return p == "/metrics" || p == "/healthz" || p == "/readyz"
	}

	stack := middleware.Chain(
		mux,
		func(next http.Handler) http.Handler { return observability.MetricsMiddleware(next) },
		middleware.RequestID(),
		middleware.CORS(cfg.CORS),
		middleware.Skip(middleware.APIKey(cfg.Auth.APIKeys), skipAuth),
		middleware.Skip(middleware.JWT(cfg.Auth.JWT), skipAuth),
		middleware.Skip(middleware.RateLimit(cfg.Rate), skipAuth),
	)

	return &http.Server{
		Addr:              cfg.HTTP.Address,
		Handler:           observability.Middleware(stack),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
	}
}
