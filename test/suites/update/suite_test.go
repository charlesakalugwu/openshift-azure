//+build e2e

package update

import (
	"flag"
	"fmt"
	"github.com/openshift/openshift-azure/test/util/client/azure"
	"github.com/openshift/openshift-azure/test/util/client/cluster"
	"github.com/openshift/openshift-azure/test/util/client/kubernetes"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	c = kubernetes.NewClient(*kubeconfig, *artifactDir)
	a = azure.NewClient()

	pluginConfig = cluster.NewPluginConfig()
	ctx = a.NewContextWithAzureCredentials()
	logger = a.NewLogger(*logLevel)
	logger.Debugf("manifest path: %s", *manifest)
	logger.Debugf("config blob path: %s", *configBlob)
	logger.Debugf("artifacts path: %s", *artifactDir)
})

func TestE2eRP(t *testing.T) {
	flag.Parse()
	fmt.Printf("E2E Resource Provider tests starting, git commit %s\n", gitCommit)
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Resource Provider Suite")
}
