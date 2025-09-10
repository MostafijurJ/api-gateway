package middleware

import (
	"context"
	"net/http"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rid := time.Now().UTC().Format("20060102T150405.000000000Z07:00")
			w.Header().Set("X-Request-ID", rid)
			ctx := context.WithValue(r.Context(), requestIDKey, rid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
