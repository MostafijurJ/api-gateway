package proxy

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"api-gateway/config"
	"api-gateway/internal/observability"
)

func NewReverseProxy(cfg config.UpstreamConfig) (http.Handler, error) {
	upstreamURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	rp := httputil.NewSingleHostReverseProxy(upstreamURL)
	rp.Transport = transport
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		observability.Logger().Error().Err(err).Str("path", r.URL.Path).Msg("upstream error")
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
	}
	rp.ModifyResponse = func(resp *http.Response) error {
		return nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), cfg.Timeout)
		defer cancel()
		r = r.WithContext(ctx)
		rp.ServeHTTP(w, r)
	}), nil
}
