//+build e2e

package cluster

import (
	"io/ioutil"

	"github.com/ghodss/yaml"

	"github.com/openshift/openshift-azure/pkg/api"
	v20180930preview "github.com/openshift/openshift-azure/pkg/api/2018-09-30-preview/api"
	"github.com/openshift/openshift-azure/pkg/config"
	"github.com/openshift/openshift-azure/pkg/fakerp"
)

// ParseExternalConfig parses an external manifest located at path and returns an external OpenshiftManagedCluster struct
func ParseExternalConfig(path string) (*v20180930preview.OpenShiftManagedCluster, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cs *v20180930preview.OpenShiftManagedCluster
	if err := yaml.Unmarshal(b, &cs); err != nil {
		return nil, err
	}
	return cs, nil
}

// SaveConfig marshall an internal OpenshiftManagedCluster struct to yaml and writes it to a file at path
func SaveConfig(config *api.OpenShiftManagedCluster, path string) error {
	if path == "" {
		path = "."
	}
	b, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, b, 0666)
	return err
}

// DeleteSecrets deletes all non-ca certificates and secrets from an internal OpenshiftManagedCluster config blob
func DeleteSecrets(config *api.OpenShiftManagedCluster) *api.OpenShiftManagedCluster {
	configCopy := config.DeepCopy()
	if configCopy.Config == nil {
		return configCopy
	}
	// remove non-ca certificates and keys
	ca := configCopy.Config.Certificates.Ca
	etcd := configCopy.Config.Certificates.EtcdCa
	frontproxy := configCopy.Config.Certificates.FrontProxyCa
	servicecatalog := configCopy.Config.Certificates.ServiceCatalogCa
	servicesigning := configCopy.Config.Certificates.ServiceSigningCa
	configCopy.Config.Certificates = api.CertificateConfig{}
	configCopy.Config.Certificates.Ca = ca
	configCopy.Config.Certificates.EtcdCa = etcd
	configCopy.Config.Certificates.FrontProxyCa = frontproxy
	configCopy.Config.Certificates.ServiceCatalogCa = servicecatalog
	configCopy.Config.Certificates.ServiceSigningCa = servicesigning

	// remove secrets
	configCopy.Config.SSHKey = nil
	configCopy.Config.RegistryHTTPSecret = nil
	configCopy.Config.RegistryConsoleOAuthSecret = ""
	configCopy.Config.ConsoleOAuthSecret = ""
	configCopy.Config.AlertManagerProxySessionSecret = nil
	configCopy.Config.AlertsProxySessionSecret = nil
	configCopy.Config.PrometheusProxySessionSecret = nil
	configCopy.Config.SessionSecretAuth = nil
	configCopy.Config.SessionSecretEnc = nil
	return configCopy
}

// GenerateInternalConfig calls generate on a OpenshiftManagedCluster
func GenerateInternalConfig(internal *api.OpenShiftManagedCluster) error {
	configGen := config.NewSimpleGenerator(fakerp.NewPluginConfig())
	return configGen.Generate(internal)
}
