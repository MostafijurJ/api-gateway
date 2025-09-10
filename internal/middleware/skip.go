package middleware

import "net/http"

// Skip wraps a middleware and skips it when predicate returns true.
func Skip(mw Middleware, predicate func(*http.Request) bool) Middleware {
	return func(next http.Handler) http.Handler {
		wrapped := mw(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if predicate(r) {
				next.ServeHTTP(w, r)
				return
			}
			wrapped.ServeHTTP(w, r)
		})
	}
}
