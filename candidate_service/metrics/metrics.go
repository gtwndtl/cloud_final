package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
    CandidatesTotal = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "candidates_total",
        Help: "Total number of candidates",
    })

    CandidatesCreateTotal = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "candidates_create_total",
        Help: "Total number of candidates created",
    })

    CandidatesUpdateTotal = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "candidates_update_total",
        Help: "Total number of candidates updated",
    })

    CandidatesDeleteTotal = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "candidates_delete_total",
        Help: "Total number of candidates deleted",
    })
)

func RegisterMetrics() {
    prometheus.MustRegister(
        CandidatesTotal,
        CandidatesCreateTotal,
        CandidatesUpdateTotal,
        CandidatesDeleteTotal,
    )
}
