package fakerp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/openshift/openshift-azure/pkg/api"
	"github.com/openshift/openshift-azure/pkg/plugin"
)

// handleClusterStatus handles admin requests for the status of control plane pods
func (s *Server) handleClusterStatus(w http.ResponseWriter, req *http.Request) {
	cs := s.read()
	if cs == nil {
		s.internalError(w, "Failed to read the internal config")
		return
	}

	config, err := GetPluginConfig()
	if err != nil {
		s.internalError(w, fmt.Sprintf("Failed to configure plugin: %v", err))
		return
	}
	p, errs := plugin.NewPlugin(s.log, config)
	if len(errs) > 0 {
		s.internalError(w, fmt.Sprintf("Failed to configure plugin: %v", err))
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, api.ContextKeyClientID, cs.Properties.ServicePrincipalProfile.ClientID)
	ctx = context.WithValue(ctx, api.ContextKeyClientSecret, cs.Properties.ServicePrincipalProfile.Secret)
	ctx = context.WithValue(ctx, api.ContextKeyTenantID, cs.Properties.AzProfile.TenantID)

	status, err := p.ClusterStatus(ctx, cs)
	if err != nil {
		s.internalError(w, fmt.Sprintf("Failed to fetch cluster status: %v", err))
		return
	}

	w.Write(status)
	s.log.Info("fetched cluster status")
}
