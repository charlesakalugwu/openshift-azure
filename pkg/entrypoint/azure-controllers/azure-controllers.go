package azurecontrollers

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"github.com/openshift/openshift-azure/pkg/api"
	"github.com/openshift/openshift-azure/pkg/cluster"
	"github.com/openshift/openshift-azure/pkg/controllers/customeradmin"
	"github.com/openshift/openshift-azure/pkg/metrics"
	"github.com/openshift/openshift-azure/pkg/util/cloudprovider"
	"github.com/openshift/openshift-azure/pkg/util/configblob"
	"github.com/openshift/openshift-azure/pkg/util/log"
)

func init() {
	prometheus.MustRegister(metrics.AzureControllersInfoGauge)
	prometheus.MustRegister(metrics.AzureControllersErrorsCounter)
	prometheus.MustRegister(metrics.AzureControllersDurationSummary)
	prometheus.MustRegister(metrics.AzureControllersInFlightGauge)
	prometheus.MustRegister(metrics.AzureControllersLastExecutedGauge)
}

func start(cfg *cmdConfig) error {
	ctx := context.Background()
	logrus.SetLevel(log.SanitizeLogLevel(cfg.LogLevel))
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log := logrus.NewEntry(logrus.StandardLogger())

	cpc, err := cloudprovider.Load("_data/_out/azure.conf")
	if err != nil {
		return err
	}

	bsc, err := configblob.GetService(ctx, log, cpc)
	if err != nil {
		return err
	}

	c := bsc.GetContainerReference(cluster.ConfigContainerName)
	blob := c.GetBlobReference(cluster.SyncBlobName)

	log.Print("reading config")
	rc, err := blob.Get(nil)
	if err != nil {
		return err
	}
	defer rc.Close()

	var cs *api.OpenShiftManagedCluster
	err = json.NewDecoder(rc).Decode(&cs)
	if err != nil {
		return err
	}

	log.Print("azure-controller pod starting")
	metrics.AzureControllersInfoGauge.With(prometheus.Labels{
		"name":           "",
		"image":          cs.Config.Images.AzureControllers,
		"period_seconds": "",
	}).Set(1)

	// TODO(charlesakalugwu): Use controller-runtime's metrics exposition when
	//  we are able to update controller-runtime to any versions > 0.1.0
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.httpPort))
	if err != nil {
		return err
	}

	mux := &http.ServeMux{}
	mux.Handle("/healthz/ready", http.HandlerFunc(readyHandler))
	mux.Handle(cfg.metricsEndpoint, promhttp.Handler())

	go http.Serve(l, mux)

	managerConfig, err := config.GetConfig()
	if err != nil {
		return err
	}

	m, err := manager.New(managerConfig, manager.Options{})
	if err != nil {
		return err
	}

	stopCh := signals.SetupSignalHandler()

	if err := customeradmin.AddToManager(ctx, log, m, stopCh); err != nil {
		return err
	}

	log.Print("starting manager")
	return m.Start(stopCh)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
