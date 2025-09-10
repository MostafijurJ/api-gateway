package middleware

import (
	"net/http"
	"strings"

	"api-gateway/config"
)

func APIKey(keys []string) Middleware {
	if len(keys) == 0 {
		return func(next http.Handler) http.Handler { return next }
	}
	set := map[string]struct{}{}
	for _, k := range keys {
		set[k] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" {
				key = r.URL.Query().Get("api_key")
			}
			if _, ok := set[key]; !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func JWT(cfg config.JWTConfig) Middleware {
	if !cfg.Enabled || cfg.Secret == "" {
		return func(next http.Handler) http.Handler { return next }
	}
	// Minimal placeholder: Verify presence of Bearer token. Replace with real JWT validation if needed.
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			// TODO: validate JWT signature/claims with cfg.Secret, cfg.Issuer, cfg.Audience.
			next.ServeHTTP(w, r)
		})
	}
}
