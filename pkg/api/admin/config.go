package admin

import (
	"crypto/x509"

	uuid "github.com/satori/go.uuid"
)

// Config holds the cluster admin config structure
type Config struct {
	// SecurityPatchPackages defines a list of rpm packages that fix security issues
	SecurityPatchPackages *[]string `json:"securityPatchPackages,omitempty"`

	// PluginVersion defines release version of the plugin used to build the cluster
	PluginVersion *string `json:"pluginVersion,omitempty"`

	// ComponentLogLevel specifies the log levels for the various openshift components
	ComponentLogLevel *ComponentLogLevel `json:"componentLogLevel,omitempty"`

	// configuration of VMs in ARM template
	ImageOffer     *string `json:"imageOffer,omitempty"`
	ImagePublisher *string `json:"imagePublisher,omitempty"`
	ImageSKU       *string `json:"imageSku,omitempty"`
	ImageVersion   *string `json:"imageVersion,omitempty"`

	// SSH to system nodes allowed IP ranges
	SSHSourceAddressPrefixes *[]string `json:"sshSourceAddressPrefixes,omitempty"`

	// configuration of other ARM resources
	ConfigStorageAccount    *string `json:"configStorageAccount,omitempty"`
	RegistryStorageAccount  *string `json:"registryStorageAccount,omitempty"`
	AzureFileStorageAccount *string `json:"azureFileStorageAccount,omitempty"`

	Certificates *CertificateConfig `json:"certificates,omitempty"`
	Images       *ImageConfig       `json:"images,omitempty"`

	// misc infra configurables
	ServiceCatalogClusterID *uuid.UUID `json:"serviceCatalogClusterId,omitempty"`

	// Geneva Metrics System (MDM) sector used for logging
	GenevaLoggingSector *string `json:"genevaLoggingSector,omitempty"`
	// Geneva Metrics System (MDM) logging account
	GenevaLoggingAccount *string `json:"genevaLoggingAccount,omitempty"`
	// Geneva Metrics System (MDM) logging namespace
	GenevaLoggingNamespace *string `json:"genevaLoggingNamespace,omitempty"`
	// Geneva Metrics System (MDM) logging control plane parameters
	GenevaLoggingControlPlaneAccount     *string `json:"genevaLoggingControlPlaneAccount,omitempty"`
	GenevaLoggingControlPlaneEnvironment *string `json:"genevaLoggingControlPlaneEnvironment,omitempty"`
	GenevaLoggingControlPlaneRegion      *string `json:"genevaLoggingControlPlaneRegion,omitempty"`
	// Geneva Metrics System (MDM) account name for metrics
	GenevaMetricsAccount *string `json:"genevaMetricsAccount,omitempty"`
	// Geneva Metrics System (MDM) endpoint for metrics
	GenevaMetricsEndpoint *string `json:"genevaMetricsEndpoint,omitempty"`
}

// ComponentLogLevel represents the log levels for the various components of a
// cluster
type ComponentLogLevel struct {
	APIServer         *int `json:"apiServer,omitempty"`
	ControllerManager *int `json:"controllerManager,omitempty"`
	Node              *int `json:"node,omitempty"`
}

// ImageConfig contains all images for the pods
type ImageConfig struct {
	// Format of the pull spec that is going to be
	// used in the cluster.
	Format *string `json:"format,omitempty"`

	ClusterMonitoringOperator *string `json:"clusterMonitoringOperator,omitempty"`
	AzureControllers          *string `json:"azureControllers,omitempty"`
	PrometheusOperator        *string `json:"prometheusOperator,omitempty"`
	Prometheus                *string `json:"prometheus,omitempty"`
	PrometheusConfigReloader  *string `json:"prometheusConfigReloader,omitempty"`
	ConfigReloader            *string `json:"configReloader,omitempty"`
	AlertManager              *string `json:"alertManager,omitempty"`
	NodeExporter              *string `json:"nodeExporter,omitempty"`
	Grafana                   *string `json:"grafana,omitempty"`
	KubeStateMetrics          *string `json:"kubeStateMetrics,omitempty"`
	KubeRbacProxy             *string `json:"kubeRbacProxy,omitempty"`
	OAuthProxy                *string `json:"oAuthProxy,omitempty"`

	MasterEtcd            *string `json:"masterEtcd,omitempty"`
	ControlPlane          *string `json:"controlPlane,omitempty"`
	Node                  *string `json:"node,omitempty"`
	ServiceCatalog        *string `json:"serviceCatalog,omitempty"`
	Sync                  *string `json:"sync,omitempty"`
	Startup               *string `json:"startup,omitempty"`
	TemplateServiceBroker *string `json:"templateServiceBroker,omitempty"`
	TLSProxy              *string `json:"tlsProxy,omitempty"`
	Registry              *string `json:"registry,omitempty"`
	Router                *string `json:"router,omitempty"`
	RegistryConsole       *string `json:"registryConsole,omitempty"`
	AnsibleServiceBroker  *string `json:"ansibleServiceBroker,omitempty"`
	WebConsole            *string `json:"webConsole,omitempty"`
	Console               *string `json:"console,omitempty"`
	EtcdBackup            *string `json:"etcdBackup,omitempty"`
	Httpd                 *string `json:"httpd,omitempty"`
	Canary                *string `json:"canary,omitempty"`

	// Geneva integration images
	GenevaLogging *string `json:"genevaLogging,omitempty"`
	GenevaTDAgent *string `json:"genevaTDAgent,omitempty"`
	GenevaStatsd  *string `json:"genevaStatsd,omitempty"`
	MetricsBridge *string `json:"metricsBridge,omitempty"`

	MonitorAgent *string `json:"monitorAgent,omitempty"`
}

// CertificateConfig contains all certificate configuration for the cluster.
type CertificateConfig struct {
	// CAs
	EtcdCa           *Certificate `json:"etcdCa,omitempty"`
	Ca               *Certificate `json:"ca,omitempty"`
	FrontProxyCa     *Certificate `json:"frontProxyCa,omitempty"`
	ServiceSigningCa *Certificate `json:"serviceSigningCa,omitempty"`
	ServiceCatalogCa *Certificate `json:"serviceCatalogCa,omitempty"`

	// etcd certificates
	EtcdServer *Certificate `json:"etcdServer,omitempty"`
	EtcdPeer   *Certificate `json:"etcdPeer,omitempty"`
	EtcdClient *Certificate `json:"etcdClient,omitempty"`

	// control plane certificates
	MasterServer *Certificate `json:"masterServer,omitempty"`
	// external web facing certificates must contain
	// all certificate chain
	// TODO: Move all certificates to be slice
	OpenShiftConsole     *CertificateChain `json:"openShiftConsole,omitempty"`
	Admin                *Certificate      `json:"admin,omitempty"`
	AggregatorFrontProxy *Certificate      `json:"aggregatorFrontProxy,omitempty"`
	MasterKubeletClient  *Certificate      `json:"masterKubeletClient,omitempty"`
	MasterProxyClient    *Certificate      `json:"masterProxyClient,omitempty"`
	OpenShiftMaster      *Certificate      `json:"openShiftMaster,omitempty"`
	NodeBootstrap        *Certificate      `json:"nodeBootstrap,omitempty"`
	SDN                  *Certificate      `json:"sdn,omitempty"`

	// infra certificates
	Registry             *Certificate      `json:"registry,omitempty"`
	RegistryConsole      *Certificate      `json:"registryConsole,omitempty"`
	Router               *CertificateChain `json:"router,omitempty"`
	ServiceCatalogServer *Certificate      `json:"serviceCatalogServer,omitempty"`

	// misc certificates
	BlackBoxMonitor *Certificate `json:"blackBoxMonitor,omitempty"`

	// geneva integration certificates
	GenevaLogging *Certificate `json:"genevaLogging,omitempty"`
	GenevaMetrics *Certificate `json:"genevaMetrics,omitempty"`

	// red hat cdn client certificates
	PackageRepository *Certificate `json:"packageRepository,omitempty"`
}

// Certificate is an x509 certificate.
type Certificate struct {
	Cert *x509.Certificate `json:"cert,omitempty"`
}

type CertificateChain struct {
	Certs []*x509.Certificate `json:"certs,omitempty"`
}
