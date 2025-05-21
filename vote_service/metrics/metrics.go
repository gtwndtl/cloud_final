package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Metrics เกี่ยวกับ Vote
	VotesCreatedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "votes_created_total",
		Help: "Total number of votes created",
	})

	VotesDeletedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "votes_deleted_total",
		Help: "Total number of votes deleted",
	})

	// Metrics HTTP Request
	HTTPRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		},
		[]string{"method", "path"},
	)

	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed",
		},
		[]string{"method", "path", "status"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(
		VotesCreatedTotal,
		VotesDeletedTotal,
		HTTPRequestDurationSeconds,
		HTTPRequestsTotal,
	)
}
