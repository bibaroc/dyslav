package server

import (
	"net/http"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

func MakeMetrics() Metrics {
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: "dyslav",
		Subsystem: "user_service",
		Name:      "request_count",
		Help:      "Total number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(prometheus.SummaryOpts{
		Namespace: "dyslav",
		Subsystem: "user_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	return Metrics{
		requestCount:   requestCount,
		requestLatency: requestLatency,
	}
}
