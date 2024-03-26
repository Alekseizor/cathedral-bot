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

var panicMetrics = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "cathedral_bot",
	Subsystem: "request",
	Name:      "panic",
	Help:      "The number of panics",
})

var handlerErrorMetrics = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "cathedral_bot",
	Subsystem: "request",
	Name:      "handler_error",
	Help:      "Counts the number of errors that occurred on the handler for each of the states",
}, []string{"state"})

var sendErrorMetrics = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "cathedral_bot",
	Subsystem: "request",
	Name:      "send_error",
	Help:      "The number of errors when sending messages for each of the states",
}, []string{"state"})

// Percentiles - считает количество запросов по каждому из стейтов,
// а также перцентили для них
func (m *Metrics) Percentiles(d time.Duration, state string) {
	requestMetrics.WithLabelValues(state).Observe(d.Seconds())
}

// Panic - считает количество паник
func (m *Metrics) Panic() {
	panicMetrics.Inc()
}

// HandlerError - считает количество ошибок, которые произошли на хендлере,
// для каждого из стейтов
func (m *Metrics) HandlerError(state string) {
	handlerErrorMetrics.WithLabelValues(state).Inc()
}

// SendError - считает количество ошибок, которые произошли при отправке,
// для каждого из стейтов
func (m *Metrics) SendError(state string) {
	sendErrorMetrics.WithLabelValues(state).Inc()
}
