package config

import (
	"reflect"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"

	"github.com/openshift/openshift-azure/pkg/api"
	"github.com/openshift/openshift-azure/pkg/util/fixtures"
)

func TestGenerate(t *testing.T) {
	tests := map[string]struct {
		cs *api.OpenShiftManagedCluster
	}{
		"test generate new": {
			cs: fixtures.NewTestOpenShiftCluster(),
		},
		//TODO "test generate doesn't overwrite": {},
	}
	var cg simpleGenerator
	for name, test := range tests {
		err := cg.Generate(test.cs)
		if err != nil {
			t.Errorf("%s received generation error %v", name, err)
			continue
		}
		testRequiredFields(test.cs, cg.pluginConfig, t)
		// check mutation
	}
}

func testRequiredFields(cs *api.OpenShiftManagedCluster, pc api.PluginConfig, t *testing.T) {
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
	assert(c.Images.LogBridge != "", "logbridge image")
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

	assert(c.ServiceAccountKey != nil, "service account key")

	if pc.TestConfig.RunningUnderTest {
		assert(c.RunningUnderTest == true, "running under test")
		assert(len(c.HtPasswd) != 0, "htpassword")
		assert(len(c.CustomerAdminPasswd) != 0, "customer-cluster-admin password")
		assert(len(c.CustomerReaderPasswd) != 0, "customer-cluster-reader password")
		assert(len(c.EndUserPasswd) != 0, "end user password")
	}
	assert(c.SSHKey != nil, "ssh key")

	assert(len(c.RegistryStorageAccount) != 0, "registry storage account")
	assert(len(c.RegistryConsoleOAuthSecret) != 0, "registry console oauth secret")
	assert(len(c.RouterStatsPassword) != 0, "router stats password")
	assert(len(c.LoggingWorkspace) != 0, "logging workspace")
	assert(len(c.LoggingLocation) != 0, "logging location")

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
	assertCert(c.Certificates.OpenshiftConsole, "OpenshiftConsole")
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

func TestGenerateServiceAccountKey(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	var cg simpleGenerator
	oc := fixtures.NewTestOpenShiftCluster()

	Expect(oc).ToNot(BeNil())
	oc.Config.ServiceAccountKey = nil
	err := cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(oc.Config.ServiceAccountKey).NotTo(BeNil())
	key := oc.Config.ServiceAccountKey
	err = cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(reflect.DeepEqual(oc.Config.ServiceAccountKey, key)).To(BeTrue())
}

func TestGenerateWhileRunningUnderTest(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	var cg simpleGenerator
	oc := fixtures.NewTestOpenShiftCluster()

	Expect(oc).ToNot(BeNil())
	oc.Config.RunningUnderTest = true
	oc.Config.HtPasswd = []byte("Y3VzdG9tZXItY2x1c3Rlci1hZG1pbjokMmEkMTAkdG1hQ1RZQVdYSVFtQ1VVMzFEOVpVZUlFSjQ0bmxENmdYeEFicTZGekUwOGo2S1Z2b0wub0sKY3VzdG9tZXItY2x1c3Rlci1yZWFkZXI6JDJhJDEwJGV0SnhmVG5XSWYuSUMyUGF4SThZbmUxbXFhUE5wcXlEcGdCWHFlaENERlQ4NGFSRHZuenAuCmVuZHVzZXI6JDJhJDEwJGRpY3BhU29aTjZObGp3N1IuOElNOHVIVDRDVmdhZ2djTVd0U1ZtWW1JWWkwVlJ5aS5FdTdxCg==")
	oc.Config.CustomerAdminPasswd = ""
	oc.Config.CustomerReaderPasswd = ""
	oc.Config.EndUserPasswd = ""
	err := cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(oc.Config.HtPasswd).NotTo(BeNil())
	Expect(len(oc.Config.HtPasswd)).To(BeNumerically(">=", 0))
	htpasswd := oc.Config.HtPasswd
	err = cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(reflect.DeepEqual(oc.Config.HtPasswd, htpasswd)).To(BeTrue())
}

func TestGenerateSSHKey(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	var cg simpleGenerator
	oc := fixtures.NewTestOpenShiftCluster()

	Expect(oc).ToNot(BeNil())
	oc.Config.SSHKey = nil
	err := cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(oc.Config.SSHKey).NotTo(BeNil())
	key := oc.Config.SSHKey
	err = cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(reflect.DeepEqual(oc.Config.SSHKey, key)).To(BeTrue())
}

func TestGenerateRegistryStorageAccount(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	var cg simpleGenerator
	oc := fixtures.NewTestOpenShiftCluster()

	Expect(oc).ToNot(BeNil())
	oc.Config.RegistryStorageAccount = ""
	err := cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(oc.Config.RegistryStorageAccount).NotTo(BeEmpty())
	name := oc.Config.RegistryStorageAccount
	err = cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(reflect.DeepEqual(oc.Config.RegistryStorageAccount, name)).To(BeTrue())
}

func TestGenerateConfigStorageAccount(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	var cg simpleGenerator
	oc := fixtures.NewTestOpenShiftCluster()

	Expect(oc).ToNot(BeNil())
	oc.Config.ConfigStorageAccount = ""
	err := cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(oc.Config.ConfigStorageAccount).NotTo(BeEmpty())
	name := oc.Config.ConfigStorageAccount
	err = cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(reflect.DeepEqual(oc.Config.ConfigStorageAccount, name)).To(BeTrue())
}

func TestGenerateLoggingWorkspace(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	var cg simpleGenerator
	oc := fixtures.NewTestOpenShiftCluster()

	Expect(oc).ToNot(BeNil())
	oc.Config.LoggingWorkspace = ""
	err := cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(oc.Config.LoggingWorkspace).NotTo(BeEmpty())
	name := oc.Config.LoggingWorkspace
	err = cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(reflect.DeepEqual(oc.Config.LoggingWorkspace, name)).To(BeTrue())
}

func TestGenerateLoggingLocation(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	var cg simpleGenerator
	oc := fixtures.NewTestOpenShiftCluster()

	Expect(oc).ToNot(BeNil())
	oc.Config.LoggingLocation = ""
	err := cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(oc.Config.LoggingLocation).NotTo(BeEmpty())
	name := oc.Config.LoggingLocation
	err = cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(reflect.DeepEqual(oc.Config.LoggingLocation, name)).To(BeTrue())
}

func TestGenerateRegistryConsoleOAuthSecret(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	var cg simpleGenerator
	oc := fixtures.NewTestOpenShiftCluster()

	Expect(oc).ToNot(BeNil())
	oc.Config.RegistryConsoleOAuthSecret = ""
	err := cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(oc.Config.RegistryConsoleOAuthSecret).NotTo(BeEmpty())
	Expect(strings.HasPrefix(oc.Config.RegistryConsoleOAuthSecret, "user")).To(BeTrue())
	Expect(len(oc.Config.RegistryConsoleOAuthSecret)).To(Equal(68))
	secret := oc.Config.RegistryConsoleOAuthSecret
	err = cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(reflect.DeepEqual(oc.Config.RegistryConsoleOAuthSecret, secret)).To(BeTrue())
}

func TestGenerateConsoleOAuthSecret(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	var cg simpleGenerator
	oc := fixtures.NewTestOpenShiftCluster()

	Expect(oc).ToNot(BeNil())
	oc.Config.ConsoleOAuthSecret = ""
	err := cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(oc.Config.ConsoleOAuthSecret).NotTo(BeEmpty())
	Expect(len(oc.Config.ConsoleOAuthSecret)).To(Equal(64))
	secret := oc.Config.ConsoleOAuthSecret
	err = cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(reflect.DeepEqual(oc.Config.ConsoleOAuthSecret, secret)).To(BeTrue())
}

func TestGenerateRouterStatsPassword(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	var cg simpleGenerator
	oc := fixtures.NewTestOpenShiftCluster()

	Expect(oc).ToNot(BeNil())
	oc.Config.RouterStatsPassword = ""
	err := cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(oc.Config.RouterStatsPassword).NotTo(BeEmpty())
	Expect(len(oc.Config.RouterStatsPassword)).To(Equal(10))
	pass := oc.Config.RouterStatsPassword
	err = cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(reflect.DeepEqual(oc.Config.RouterStatsPassword, pass)).To(BeTrue())
}

func TestGenerateServiceCatalogClusterID(t *testing.T) {
	RegisterTestingT(t)
	t.Parallel()

	var cg simpleGenerator
	oc := fixtures.NewTestOpenShiftCluster()

	Expect(oc).ToNot(BeNil())
	oc.Config.ServiceCatalogClusterID = uuid.Nil
	err := cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(uuid.Equal(oc.Config.ServiceCatalogClusterID, uuid.Nil)).To(BeFalse())
	uuid := oc.Config.ServiceCatalogClusterID
	err = cg.Generate(oc)
	Expect(err).To(BeNil())
	Expect(reflect.DeepEqual(oc.Config.ServiceCatalogClusterID, uuid)).To(BeTrue())
}
