package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	once    sync.Once
	metrics *PromMetrics
)

type PromMetrics struct {
	ConsumeCounter    *prometheus.CounterVec
	DB_duration       *prometheus.HistogramVec
	Response_duration *prometheus.HistogramVec
}

func BuildMetrics() {
	once.Do(func() {
		NewPromMetrics()
	})
}
func NewPromMetrics() {
	//p.msgCounter.Inc()
	consumeCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "kata_test_server",
		Name:      "GetRates_counter",
		Help:      "GetRates_counter",
	}, []string{"GetRates_handler_counter"})

	response_duration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "kata_test_server",
		Name:      "response_latency_seconds",
		Help:      "response_latency",
		Buckets:   []float64{0.1, 0.3, 0.5, 1}, //from one hundred milliseconds to a second
	}, []string{"GetRates_handler_duration"})

	db_duration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "kata_test_server",
		Name:      "db_latency_seconds",
		Help:      "db_latency",
		Buckets:   []float64{0.1, 0.3, 0.5, 1}, //from one hundred milliseconds to a second
	}, []string{"GetRates_db_duration"})

	metrics = &PromMetrics{
		ConsumeCounter:    consumeCounter,
		DB_duration:       db_duration,
		Response_duration: response_duration,
	}

}

func Metrics() *PromMetrics {
	if metrics == nil {
		BuildMetrics()
	}
	return metrics
}
