package config

import (
	"reflect"
	"testing"

	"github.com/satori/go.uuid"

	"github.com/openshift/openshift-azure/pkg/api"
	pluginapi "github.com/openshift/openshift-azure/pkg/api/plugin/api"
	"github.com/openshift/openshift-azure/test/util/populate"
)

func TestGenerate(t *testing.T) {
	cs := &api.OpenShiftManagedCluster{
		Properties: api.Properties{
			OpenShiftVersion: "v3.11",
			RouterProfiles: []api.RouterProfile{
				{},
			},
		},
	}
	pc := api.PluginConfig{
		TestConfig: api.TestConfig{
			RunningUnderTest: true,
		},
	}

	prepare := func(v reflect.Value) {}
	var template *pluginapi.Config
	populate.Walk(&template, prepare)

	cg := simpleGenerator{pluginConfig: pc}
	err := cg.Generate(cs, template)
	if err != nil {
		t.Error(err)
	}

	testRequiredFields(cs, t)
}

func testRequiredFields(cs *api.OpenShiftManagedCluster, t *testing.T) {
	assert := func(c bool, name string) {
		if !c {
			t.Errorf("missing %s", name)
		}
	}
	assertCert := func(c api.CertKeyPair, name string) {
		assert(c.Key != nil, name+" key")
		assert(c.Cert != nil, name+" cert")
	}

	c := cs.Config

	assert(c.ImagePublisher != "", "image publisher")
	assert(c.ImageOffer != "", "image offer")
	assert(c.ImageVersion != "", "image version")

	assert(c.Images.Format != "", "image config format")
	assert(c.Images.ControlPlane != "", "control plane image")
	assert(c.Images.Node != "", "node image")
	assert(c.Images.ServiceCatalog != "", "service catalog image")
	assert(c.Images.AnsibleServiceBroker != "", "ansible service broker image")
	assert(c.Images.TemplateServiceBroker != "", "template service broker image")
	assert(c.Images.Registry != "", "registry image")
	assert(c.Images.Router != "", "router image")
	assert(c.Images.WebConsole != "", "web console image")
	assert(c.Images.Console != "", "console image")
	assert(c.Images.MasterEtcd != "", "master etcd image")
	assert(c.Images.RegistryConsole != "", "registry console image")
	assert(c.Images.Sync != "", "sync image")
	assert(c.Images.EtcdBackup != "", "etcdbackup image")
	assert(c.Images.ClusterMonitoringOperator != "", "cluster monitoring operator image")
	assert(c.Images.PrometheusOperatorBase != "", "cluster monitoring operator image")
	assert(c.Images.PrometheusConfigReloaderBase != "", "prometheus config reloader base image")
	assert(c.Images.ConfigReloaderBase != "", "config reloader base image")
	assert(c.Images.PrometheusBase != "", "prometheus base image")
	assert(c.Images.AlertManagerBase != "", "alertmanager base image")
	assert(c.Images.NodeExporterBase != "", "node exporter base image")
	assert(c.Images.GrafanaBase != "", "grafana base image")
	assert(c.Images.KubeStateMetricsBase != "", "kube state metrics base image")
	assert(c.Images.KubeRbacProxyBase != "", "kube rbac proxy base image")
	assert(c.Images.OAuthProxyBase != "", "oauth proxy base image")
	assert(c.Images.GenevaLogging != "", "azure logging image")
	assert(c.Images.GenevaTDAgent != "", "azure TDAgent image")
	assert(c.Images.MetricsBridge != "", "metrics-bridge image")

	assert(c.ServiceAccountKey != nil, "service account key")

	assert(len(c.HtPasswd) != 0, "htpassword")
	assert(len(c.CustomerAdminPasswd) != 0, "customer-cluster-admin password")
	assert(len(c.CustomerReaderPasswd) != 0, "customer-cluster-reader password")
	assert(len(c.EndUserPasswd) != 0, "end user password")

	assert(c.SSHKey != nil, "ssh key")

	assert(len(c.RegistryStorageAccount) != 0, "registry storage account")
	assert(len(c.RegistryConsoleOAuthSecret) != 0, "registry console oauth secret")
	assert(len(c.RouterStatsPassword) != 0, "router stats password")

	assert(c.ServiceCatalogClusterID != uuid.Nil, "service catalog cluster id")

	assertCert(c.Certificates.EtcdCa, "EtcdCa")
	assertCert(c.Certificates.Ca, "Ca")
	assertCert(c.Certificates.FrontProxyCa, "FrontProxyCa")
	assertCert(c.Certificates.ServiceSigningCa, "ServiceSigningCa")
	assertCert(c.Certificates.ServiceCatalogCa, "ServiceCatalogCa")
	assertCert(c.Certificates.EtcdServer, "EtcdServer")
	assertCert(c.Certificates.EtcdPeer, "EtcdPeer")
	assertCert(c.Certificates.EtcdClient, "EtcdClient")
	assertCert(c.Certificates.MasterServer, "MasterServer")
	assertCert(c.Certificates.OpenShiftConsole, "OpenShiftConsole")
	assertCert(c.Certificates.Admin, "Admin")
	assertCert(c.Certificates.AggregatorFrontProxy, "AggregatorFrontProxy")
	assertCert(c.Certificates.MasterKubeletClient, "MasterKubeletClient")
	assertCert(c.Certificates.MasterProxyClient, "MasterProxyClient")
	assertCert(c.Certificates.OpenShiftMaster, "OpenShiftMaster")
	assertCert(c.Certificates.NodeBootstrap, "NodeBootstrap")
	assertCert(c.Certificates.Registry, "Registry")
	assertCert(c.Certificates.Router, "Router")
	assertCert(c.Certificates.ServiceCatalogServer, "ServiceCatalogServer")
	assertCert(c.Certificates.ServiceCatalogAPIClient, "ServiceCatalogAPIClient")
	assertCert(c.Certificates.AzureClusterReader, "AzureClusterReader")

	assert(len(c.SessionSecretAuth) != 0, "SessionSecretAuth")
	assert(len(c.SessionSecretEnc) != 0, "SessionSecretEnc")
	assert(len(c.RegistryHTTPSecret) != 0, "RegistryHTTPSecret")
	assert(len(c.AlertManagerProxySessionSecret) != 0, "AlertManagerProxySessionSecret")
	assert(len(c.AlertsProxySessionSecret) != 0, "AlertsProxySessionSecret")
	assert(len(c.PrometheusProxySessionSecret) != 0, "PrometheusProxySessionSecret")

	assert(c.MasterKubeconfig != nil, "MasterKubeconfig")
	assert(c.AdminKubeconfig != nil, "AdminKubeconfig")
	assert(c.NodeBootstrapKubeconfig != nil, "NodeBootstrapKubeconfig")
	assert(c.AzureClusterReaderKubeconfig != nil, "AzureClusterReaderKubeconfig")
}

func TestInvalidateSecrets(t *testing.T) {
	dpc := &api.PluginConfig{}
	populate.Walk(dpc, nil)
	cs := &api.OpenShiftManagedCluster{}
	populate.Walk(cs, nil)

	g := &simpleGenerator{
		pluginConfig: *dpc,
	}
	saved := *cs
	if err := g.InvalidateSecrets(cs); err != nil {
		t.Errorf("configGenerator.InvalidateSecrets error = %v", err)
	}
	if !reflect.DeepEqual(saved.Config, cs.Config) {
		if !reflect.DeepEqual(saved.Config.Certificates.Ca, cs.Config.Certificates.Ca) {
			t.Errorf("unexpected change to ca certificates after secret invalidation")
		}
		if !reflect.DeepEqual(saved.Config.Certificates.EtcdCa, cs.Config.Certificates.EtcdCa) {
			t.Errorf("unexpected change to etcd ca certificates after secret invalidation")
		}
		if !reflect.DeepEqual(saved.Config.Certificates.FrontProxyCa, cs.Config.Certificates.FrontProxyCa) {
			t.Errorf("unexpected change to front proxy ca certificates after secret invalidation")
		}
		if !reflect.DeepEqual(saved.Config.Certificates.ServiceCatalogCa, cs.Config.Certificates.ServiceCatalogCa) {
			t.Errorf("unexpected change to service catalog ca certificates after secret invalidation")
		}
		if !reflect.DeepEqual(saved.Config.Certificates.ServiceSigningCa, cs.Config.Certificates.ServiceSigningCa) {
			t.Errorf("unexpected change to service signing ca certificates after secret invalidation")
		}
		if !reflect.DeepEqual(saved.Config.ServiceAccountKey, cs.Config.ServiceAccountKey) {
			t.Errorf("unexpected change to service service account key after secret invalidation")
		}
	} else {
		t.Errorf("expected config blob to be different after secret invalidation")
	}
}
