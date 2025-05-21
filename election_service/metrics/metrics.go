package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// จำนวน election ทั้งหมดในระบบ (gauge)
	ElectionsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "elections_total",
		Help: "Total number of elections",
	})

	// นับจำนวนการสร้าง election
	ElectionsCreateTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "elections_create_total",
		Help: "Total number of elections created",
	})

	// นับจำนวนการแก้ไข election
	ElectionsUpdateTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "elections_update_total",
		Help: "Total number of elections updated",
	})

	// นับจำนวนการลบ election
	ElectionsDeleteTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "elections_delete_total",
		Help: "Total number of elections deleted",
	})

	HTTPRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		},
		[]string{"method", "path"},
	)


	// นับจำนวน HTTP requests ที่รับมา
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
		ElectionsTotal,
		ElectionsCreateTotal,
		ElectionsUpdateTotal,
		ElectionsDeleteTotal,
		HTTPRequestDurationSeconds,
		HTTPRequestsTotal,
	)
}
