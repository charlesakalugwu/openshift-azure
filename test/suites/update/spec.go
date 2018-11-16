//+build e2e

package update

import (
	"context"
	"flag"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/openshift/openshift-azure/pkg/api"
	"github.com/openshift/openshift-azure/test/util/client/azure"
	"github.com/openshift/openshift-azure/test/util/client/cluster"
	"github.com/openshift/openshift-azure/test/util/client/kubernetes"
	"github.com/openshift/openshift-azure/test/util/scenarios/enduser"
)

var (
	c            *kubernetes.Client
	a            *azure.Client
	gitCommit    = "unknown"
	logger       *logrus.Entry
	ctx          context.Context
	pluginConfig *api.PluginConfig

	manifest    = flag.String("manifest", "../../../_data/manifest.yaml", "Path to the manifest to send to the RP")
	configBlob  = flag.String("configBlob", "../../../_data/containerservice.yaml", "Path on disk where the OpenShift internal config blob should be written")
	logLevel    = flag.String("logLevel", "Debug", "The log level to use")
	kubeconfig  = flag.String("kubeconfig", "../../../_data/_out/admin.kubeconfig", "Location of the kubeconfig")
	artifactDir = flag.String("artifact-dir", "../../../_data/_artifacts/", "Directory to place artifacts when a test fails")
)

var _ = Describe("Resource provider e2e tests [Fake] [KeyRotation]", func() {
	defer GinkgoRecover()

	It("should be possible to maintain a healthy cluster after rotating all credentials", func() {
		By("Parsing the external manifest")
		external, err := cluster.ParseExternalConfig(*manifest)
		Expect(err).NotTo(HaveOccurred())
		Expect(external).NotTo(BeNil())

		By("Parsing the internal manifest containing config blob")
		internal, err := cluster.ParseInternalConfig(*configBlob)
		Expect(err).NotTo(HaveOccurred())
		Expect(internal).NotTo(BeNil())

		By("Deleting all non-ca cluster certificates and credentials from the config blob")
		mutated := cluster.DeleteCertificates(internal)
		Expect(err).NotTo(HaveOccurred())
		Expect(mutated).NotTo(BeNil())

		By("Running generate on the modified config blob")
		err = cluster.GenerateInternalConfig(mutated, pluginConfig)
		Expect(err).NotTo(HaveOccurred())

		By("Persisting the config blob containing the new certificates and credentials")
		err = cluster.SaveConfig(mutated, *configBlob)
		Expect(err).NotTo(HaveOccurred())

		By("Calling update on the fake rp with the updated config blob")
		updated, err := cluster.UpdateCluster(ctx, external, *configBlob, logger, pluginConfig)
		Expect(err).NotTo(HaveOccurred())
		Expect(updated).NotTo(BeNil())

		By("Parsing the config blob after the update")
		internalAfterUpdate, err := cluster.ParseInternalConfig(*configBlob)
		Expect(err).NotTo(HaveOccurred())
		Expect(internalAfterUpdate).NotTo(BeNil())

		By("Verifying that the initial config blob does not match the one created after the update")
		configMatch := reflect.DeepEqual(internal.Config.Certificates, internalAfterUpdate.Config.Certificates)
		Expect(configMatch).To(BeFalse())

		By("Verifying that the mutated config blob matches the one created after the update")
		configMatch = reflect.DeepEqual(mutated.Config.Certificates, internalAfterUpdate.Config.Certificates)
		Expect(configMatch).To(BeTrue())
	})

	Context("when the key rotation is complete", func() {
		BeforeEach(func() {
			namespace := c.GenerateRandomName("e2e-test-")
			c.CreateProject(namespace)
		})

		AfterEach(func() {
			if CurrentGinkgoTestDescription().Failed {
				if err := c.DumpInfo(); err != nil {
					logrus.Warn(err)
				}
			}
			c.CleanupProject(10 * time.Minute)
		})

		It("should disallow PDB mutations", func() {
			enduser.CheckPdbMutationsDisallowed(c)
		})

		It("should deploy a template and ensure a given text is in the contents", func() {
			enduser.CheckCanDeployTemplate(c)
		})

		It("should not crud infra resources", func() {
			enduser.CheckCrudOnInfraDisallowed(c)
		})

		It("should deploy a template with persistent storage and test failure modes", func() {
			enduser.CheckCanDeployTemplateWithPV(c)
		})
	})
})
