package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	UsersTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "users_total",
		Help: "Total number of users registered",
	})

	UserSignupsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "users_signup_total",
		Help: "Total number of successful user signups",
	})

	UserSignupFailuresTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "users_signup_failures_total",
		Help: "Total number of failed user signups",
	})

	UserLoginsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "users_logins_total",
		Help: "Total number of successful user logins",
	})

	UserLoginFailuresTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "users_login_failures_total",
		Help: "Total number of failed user logins",
	})

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
		UsersTotal,
		UserSignupsTotal,
		UserSignupFailuresTotal,
		UserLoginsTotal,
		UserLoginFailuresTotal,
		HTTPRequestDurationSeconds,
		HTTPRequestsTotal,
	)
}
