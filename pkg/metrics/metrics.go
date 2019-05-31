package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// TODO(charlesakalugwu): Add unit tests for the handling of these metrics once
//  the upstream library supports it
var (
	AzureControllersInfoGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_controllers_info",
			Help: "General information about the azure controllers process.",
		},
		[]string{"name", "image", "period_seconds"},
	)

	AzureControllersErrorsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "azure_controllers_errors_total",
			Help: "Total number of errors encountered during azure controller reconciliations.",
		},
		[]string{"name"},
	)

	AzureControllersInFlightGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_controllers_executions_inflight",
			Help: "Number of azure controller reconcile executions in progress.",
		},
		[]string{"name"},
	)

	AzureControllersLastExecutedGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_controllers_last_executed",
			Help: "The last time the azure controllers were executed.",
		},
		[]string{"name"},
	)

	AzureControllersDurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "azure_controllers_duration_seconds",
			Help: "The duration of azure controller runs.",
		},
		[]string{"name"},
	)
)
