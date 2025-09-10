package router

import (
	"net/http"
	"strings"

	"api-gateway/config"
)

type Router struct {
	routes []config.Route
}

func New(routes []config.Route) *Router { return &Router{routes: routes} }

func (r *Router) Match(req *http.Request) *config.Route {
	for i := range r.routes {
		rt := &r.routes[i]
		if !strings.HasPrefix(req.URL.Path, rt.PathPrefix) {
			continue
		}
		if len(rt.Methods) > 0 {
			ok := false
			for _, m := range rt.Methods {
				if req.Method == m {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		if len(rt.Headers) > 0 {
			all := true
			for k, v := range rt.Headers {
				if req.Header.Get(k) != v {
					all = false
					break
				}
			}
			if !all {
				continue
			}
		}
		return rt
	}
	return nil
}
