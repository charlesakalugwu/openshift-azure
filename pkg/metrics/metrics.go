package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

// Registrar is used to selectively register metrics for various ARO commands
type Registrar interface {
	// RegisterSync registers metrics for sync
	RegisterSync() http.Handler

	// RegisterAzureControllers registers metrics for azure-controllers
	RegisterAzureControllers() http.Handler
}

// Collector holds all metrics used by all ARO commands as well as a metrics
// registry for selectively registering the metrics
type Collector struct {
	Registry *prometheus.Registry

	// sync pod metrics
	SyncInfoGauge         *prometheus.GaugeVec
	SyncErrorsCounter     prometheus.Counter
	SyncInFlightGauge     prometheus.Gauge
	SyncLastExecutedGauge prometheus.Gauge
	SyncDurationSummary   prometheus.Summary

	// azure-controllers metrics
	AzureControllersErrorsCounter     *prometheus.CounterVec
	AzureControllersInFlightGauge     *prometheus.GaugeVec
	AzureControllersLastExecutedGauge *prometheus.GaugeVec
	AzureControllersDurationSummary   *prometheus.SummaryVec
}

var _ Registrar = &Collector{}

// DefaultCollector returns an ARO metrics collector with its internal registry
// initialized with the go and process collectors
func DefaultCollector() *Collector {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		prometheus.NewGoCollector(),
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
	)
	return &Collector{
		Registry: registry,
	}
}

// NewCollector returns an ARO metrics collector configured with the provided
// registry
func NewCollector(registry *prometheus.Registry) *Collector {
	if registry == nil {
		return DefaultCollector()
	}
	return &Collector{
		Registry: registry,
	}
}
