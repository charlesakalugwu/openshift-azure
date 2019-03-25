package fakerp

import (
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"

	"github.com/openshift/openshift-azure/pkg/api"
	pluginapi "github.com/openshift/openshift-azure/pkg/api/plugin/api"
	"github.com/openshift/openshift-azure/pkg/util/tls"
)

func GetTestConfig() api.TestConfig {
	return api.TestConfig{
		RunningUnderTest:   os.Getenv("RUNNING_UNDER_TEST") == "true",
		ImageResourceGroup: os.Getenv("IMAGE_RESOURCEGROUP"),
		ImageResourceName:  os.Getenv("IMAGE_RESOURCENAME"),
	}
}

func GetPluginTemplate() (*pluginapi.Config, error) {
	// read template file without secrets
	data, err := ioutil.ReadFile("pluginconfig/pluginconfig-311.yaml")
	if err != nil {
		return nil, err
	}
	var template *pluginapi.Config
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, err
	}

	// enrich template with secrets
	logCert, err := readCert("secrets/logging-int.cert")
	if err != nil {
		return nil, err
	}
	logKey, err := readKey("secrets/logging-int.key")
	if err != nil {
		return nil, err
	}
	metCert, err := readCert("secrets/metrics-int.cert")
	if err != nil {
		return nil, err
	}
	metKey, err := readKey("secrets/metrics-int.key")
	if err != nil {
		return nil, err
	}
	pullSecret, err := ioutil.ReadFile("secrets/.dockerconfigjson")
	if err != nil {
		return nil, err
	}
	imagePullSecret, err := ioutil.ReadFile("secrets/system-docker-config.json")
	if err != nil {
		return nil, err
	}
	template.Certificates.GenevaLogging.Cert = logCert
	template.Certificates.GenevaLogging.Key = logKey
	template.Certificates.GenevaMetrics.Cert = metCert
	template.Certificates.GenevaMetrics.Key = metKey
	template.Images.GenevaImagePullSecret = pullSecret
	template.Images.ImagePullSecret = imagePullSecret

	return template, nil
}

func overridePluginTemplate(template *pluginapi.Config) error {
	// read plugin template override
	data, err := ioutil.ReadFile("pluginconfig/override.yaml")
	if err != nil {
		return err
	}
	var override *pluginapi.Config
	if err := yaml.Unmarshal(data, &override); err != nil {
		return err
	}

	if override.Images.Sync != "" {
		template.Images.Sync = override.Images.Sync
	}
	if override.Images.MetricsBridge != "" {
		template.Images.MetricsBridge = override.Images.MetricsBridge
	}
	if override.Images.EtcdBackup != "" {
		template.Images.EtcdBackup = override.Images.EtcdBackup
	}
	if override.Images.TLSProxy != "" {
		template.Images.TLSProxy = override.Images.TLSProxy
	}
	if override.Images.Canary != "" {
		template.Images.Canary = override.Images.Canary
	}
	if override.Images.AzureControllers != "" {
		template.Images.AzureControllers = override.Images.AzureControllers
	}
	if override.Images.Startup != "" {
		template.Images.Startup = override.Images.Startup
	}
	if override.Images.Format != "" {
		template.Images.Format = override.Images.Format
	}
	if override.ImageVersion != "" {
		template.ImageVersion = override.ImageVersion
	}
	if override.ImageOffer != "" {
		template.ImageOffer = override.ImageOffer
	}
	if override.ComponentLogLevel.APIServer >= 0 {
		template.ComponentLogLevel.APIServer = override.ComponentLogLevel.APIServer
	}
	if override.ComponentLogLevel.ControllerManager >= 0 {
		template.ComponentLogLevel.ControllerManager = override.ComponentLogLevel.ControllerManager
	}
	if override.ComponentLogLevel.Node >= 0 {
		template.ComponentLogLevel.Node = override.ComponentLogLevel.Node
	}
	return nil
}

func readCert(path string) (*x509.Certificate, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return tls.ParseCert(b)
}

func readKey(path string) (*rsa.PrivateKey, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return tls.ParsePrivateKey(b)
}
