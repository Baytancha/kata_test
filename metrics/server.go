package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MetricsServer(metrics *PromMetrics) *http.Server {

	reg := prometheus.NewRegistry() //non-global
	reg.MustRegister(metrics.ConsumeCounter)
	reg.MustRegister(metrics.Response_duration)
	reg.MustRegister(metrics.DB_duration)

	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})
	mux := http.NewServeMux()
	mux.Handle("/metrics", promHandler)
	return &http.Server{Addr: ":8182", Handler: mux}
}
