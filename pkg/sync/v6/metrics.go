package sync

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
)

var metricsRegistry = prometheus.NewRegistry()

type metrics struct {
	registry              *prometheus.Registry
	syncInfoGauge         prometheus.GaugeVec
	syncErrorsCounter     prometheus.Counter
	syncInFlightGauge     prometheus.Gauge
	syncLastExecutedGauge prometheus.Gauge
	syncDurationSummary   prometheus.Summary
}

func initMetrics() (*metrics, error) {
	m := &metrics{}
	m.registry = prometheus.NewRegistry()

	m.syncInfoGauge = *prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sync_info",
			Help: "General information about the sync process.",
		},
		[]string{"plugin_version", "image", "period_seconds"},
	)
	m.syncErrorsCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sync_errors_total",
			Help: "Total number of errors encountered during sync executions.",
		},
	)
	m.syncInFlightGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "sync_executions_inflight",
			Help: "Number of sync executions in progress.",
		},
	)
	m.syncLastExecutedGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "sync_last_executed",
			Help: "The last time a sync was executed.",
		},
	)
	m.syncDurationSummary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "sync_duration_seconds",
			Help: "The duration of sync executions.",
		},
	)
	m.registry.MustRegister(
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		m.syncInfoGauge,
		m.syncErrorsCounter,
		m.syncInFlightGauge,
		m.syncLastExecutedGauge,
		m.syncDurationSummary,
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		return nil, err
	}

	mux := &http.ServeMux{}
	//mux.Handle("/healthz/ready", http.HandlerFunc(s.ReadyHandler))
	mux.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}))

	go http.Serve(l, mux)

	return m, nil
}
