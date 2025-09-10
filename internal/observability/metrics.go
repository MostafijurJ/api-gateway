package observability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reqCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "gateway_http_requests_total", Help: "Total HTTP requests"},
		[]string{"code", "method", "path"},
	)
	reqLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "gateway_http_request_duration_seconds", Help: "Request latency", Buckets: prometheus.DefBuckets},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(reqCounter, reqLatency)
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(sr, r)
		path := r.URL.Path
		method := r.Method
		code := strconv.Itoa(sr.status)
		reqCounter.WithLabelValues(code, method, path).Inc()
		reqLatency.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
	})
}

func MetricsHandler() http.Handler { return promhttp.Handler() }
