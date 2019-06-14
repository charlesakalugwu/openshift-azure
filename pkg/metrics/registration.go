package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: Add unit tests for the handling of these metrics once
//  the upstream library supports it

// RegisterSync registers sync metrics unto the collector's registry.
func (c *Collector) RegisterSync() http.Handler {
	c.SyncInfoGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sync_info",
			Help: "General information about the sync process.",
		},
		[]string{"plugin_version", "image", "period_seconds"},
	)
	c.SyncErrorsCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sync_errors_total",
			Help: "Total number of errors encountered during sync executions.",
		},
	)
	c.SyncInFlightGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "sync_executions_inflight",
			Help: "Number of sync executions in progress.",
		},
	)
	c.SyncLastExecutedGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "sync_last_executed",
			Help: "The last time a sync was executed.",
		},
	)
	c.SyncDurationSummary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "sync_duration_seconds",
			Help: "The duration of sync executions.",
		},
	)
	c.Registry.MustRegister(
		c.SyncInfoGauge,
		c.SyncErrorsCounter,
		c.SyncInFlightGauge,
		c.SyncLastExecutedGauge,
		c.SyncDurationSummary,
	)
	return promhttp.HandlerFor(c.Registry, promhttp.HandlerOpts{})
}

// RegisterAzureControllers registers azure-controller metrics unto the
// collector's registry.
func (c *Collector) RegisterAzureControllers() http.Handler {
	c.AzureControllersErrorsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "azure_controllers_errors_total",
			Help: "Total number of errors.",
		},
		[]string{"controller"},
	)

	c.AzureControllersInFlightGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_controllers_reconciliations_inflight",
			Help: "Number of azure controller reconciliations in progress.",
		},
		[]string{"controller"},
	)

	c.AzureControllersLastExecutedGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_controllers_last_executed",
			Help: "The last time the azure controllers were run.",
		},
		[]string{"controller"},
	)

	c.AzureControllersDurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "azure_controllers_duration_seconds",
			Help: "The duration of azure controller runs.",
		},
		[]string{"controller"},
	)
	c.Registry.MustRegister(
		c.AzureControllersErrorsCounter,
		c.AzureControllersInFlightGauge,
		c.AzureControllersLastExecutedGauge,
		c.AzureControllersDurationSummary,
	)
	return promhttp.HandlerFor(c.Registry, promhttp.HandlerOpts{})
}
