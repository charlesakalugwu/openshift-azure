package fakerp

import (
	"context"
	"net/http"
	"os"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/openshift/openshift-azure/test/clients/azure"
)

var _ = Describe("Key Rotation E2E tests [KeyRotation][Fake][LongRunning]", func() {
	var (
		cli *azure.Client
	)

	BeforeEach(func() {
		var err error
		cli, err = azure.NewClientFromEnvironment(false)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should be possible to maintain a healthy cluster after rotating all credentials", func() {
		By("Reading the cluster state")
		before, err := cli.OpenShiftManagedClustersAdmin.Get(context.Background(), os.Getenv("RESOURCEGROUP"), os.Getenv("RESOURCEGROUP"))
		Expect(err).NotTo(HaveOccurred())
		Expect(before).NotTo(BeNil())

		By("Executing key rotation on the cluster.")
		update, err := cli.OpenShiftManagedClustersAdmin.RotateSecretsAndWait(context.Background(), os.Getenv("RESOURCEGROUP"), os.Getenv("RESOURCEGROUP"), before)
		Expect(err).NotTo(HaveOccurred())
		Expect(update.StatusCode).To(Equal(http.StatusOK))
		Expect(update).NotTo(BeNil())

		By("Reading the cluster state after the update")
		after, err := cli.OpenShiftManagedClustersAdmin.Get(context.Background(), os.Getenv("RESOURCEGROUP"), os.Getenv("RESOURCEGROUP"))
		Expect(err).NotTo(HaveOccurred())
		Expect(after).NotTo(BeNil())

		By("Verifying that the config blob before and after key rotation does not contain matching public certificates")
		Expect(reflect.DeepEqual(before.Config.Certificates.ServiceSigningCa.Cert, after.Config.Certificates.ServiceSigningCa.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.ServiceCatalogCa.Cert, after.Config.Certificates.ServiceCatalogCa.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.FrontProxyCa.Cert, after.Config.Certificates.FrontProxyCa.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.EtcdCa.Cert, after.Config.Certificates.EtcdCa.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.Ca.Cert, after.Config.Certificates.Ca.Cert)).To(BeFalse())

		Expect(reflect.DeepEqual(before.Config.Certificates.Admin.Cert, after.Config.Certificates.Admin.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.AggregatorFrontProxy.Cert, after.Config.Certificates.AggregatorFrontProxy.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.AzureClusterReader.Cert, after.Config.Certificates.AzureClusterReader.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.EtcdClient.Cert, after.Config.Certificates.EtcdClient.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.EtcdPeer.Cert, after.Config.Certificates.EtcdPeer.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.EtcdServer.Cert, after.Config.Certificates.EtcdServer.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.GenevaLogging.Cert, after.Config.Certificates.GenevaLogging.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.GenevaMetrics.Cert, after.Config.Certificates.GenevaMetrics.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.MasterKubeletClient.Cert, after.Config.Certificates.MasterKubeletClient.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.MasterProxyClient.Cert, after.Config.Certificates.MasterProxyClient.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.MasterServer.Cert, after.Config.Certificates.MasterServer.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.NodeBootstrap.Cert, after.Config.Certificates.NodeBootstrap.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.OpenShiftConsole.Cert, after.Config.Certificates.OpenShiftConsole.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.OpenShiftMaster.Cert, after.Config.Certificates.OpenShiftMaster.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.Registry.Cert, after.Config.Certificates.Registry.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.Router.Cert, after.Config.Certificates.Router.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.ServiceCatalogAPIClient.Cert, after.Config.Certificates.ServiceCatalogAPIClient.Cert)).To(BeFalse())
		Expect(reflect.DeepEqual(before.Config.Certificates.ServiceCatalogServer.Cert, after.Config.Certificates.ServiceCatalogServer.Cert)).To(BeFalse())

		By("Verifying that only the certificates have been updated")
		before.Config.Certificates = after.Config.Certificates
		configMatch := reflect.DeepEqual(before.Config, after.Config)
		Expect(configMatch).To(BeTrue())
	})
})
