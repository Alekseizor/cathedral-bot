package prometheus

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var requestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "cathedral_bot",
	Subsystem:  "request",
	Name:       "latency",
	Help:       "Request processing time depending on the state",
	Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.03, 0.9: 0.01, 0.99: 0.001},
}, []string{"state"})

// InformationMetrics - считает количество запросов по каждому из стейтов,
// а также перцентили для них
func (m *Metrics) InformationMetrics(d time.Duration, state string) {
	requestMetrics.WithLabelValues(state).Observe(d.Seconds())
}
