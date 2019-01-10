package fakerp

import (
	"context"
	"fmt"
	"github.com/openshift/openshift-azure/pkg/api"
	"github.com/openshift/openshift-azure/pkg/plugin"
	"net/http"
)

func (s *Server) handleRotateSecrets(w http.ResponseWriter, req *http.Request) {
	defer func() {
		// drain once we are done processing this request
		<-s.inProgress
	}()

	cs := s.read()
	if cs == nil {
		s.internalError(w, "Failed to read the internal config")
		return
	}

	ctx := context.Background()
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

	ctx = context.WithValue(ctx, api.ContextKeyClientID, cs.Properties.ServicePrincipalProfile.ClientID)
	ctx = context.WithValue(ctx, api.ContextKeyClientSecret, cs.Properties.ServicePrincipalProfile.Secret)
	ctx = context.WithValue(ctx, api.ContextKeyTenantID, cs.Properties.AzProfile.TenantID)

	deployer := GetDeployer(s.log, cs, config)
	if err := p.RecoverEtcdCluster(ctx, cs, deployer, ""); err != nil {
		s.internalError(w, fmt.Sprintf("Failed to recover cluster: %v", err))
		return
	}

	s.log.Info("recovered cluster")
}

