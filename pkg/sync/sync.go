package sync

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/openshift/openshift-azure/pkg/api"
	"github.com/openshift/openshift-azure/pkg/metrics"
	v4 "github.com/openshift/openshift-azure/pkg/sync/v4"
	v5 "github.com/openshift/openshift-azure/pkg/sync/v5"
	v6 "github.com/openshift/openshift-azure/pkg/sync/v6"
)

type Interface interface {
	Sync(ctx context.Context) error
	ReadyHandler(w http.ResponseWriter, r *http.Request)
	PrintDB() error
	Hash() ([]byte, error)
}

func New(log *logrus.Entry, cs *api.OpenShiftManagedCluster, initClients bool, metrics *metrics.Collector) (Interface, error) {
	switch cs.Config.PluginVersion {
	case "v4.2", "v4.3", "v4.4":
		return v4.New(log, cs, initClients)
	case "v5.1":
		return v5.New(log, cs, initClients)
	case "v6.0":
		return v6.New(log, cs, initClients, metrics)
	}

	return nil, fmt.Errorf("version %q not found", cs.Config.PluginVersion)
}
