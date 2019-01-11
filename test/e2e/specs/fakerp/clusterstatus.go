package fakerp

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"

	"github.com/openshift/openshift-azure/test/clients/azure"
)

var _ = Describe("Cluster Status E2E tests [ClusterStatus][Fake][EveryPR]", func() {
	var (
		cli *azure.Client
	)

	BeforeEach(func() {
		var err error
		cli, err = azure.NewClientFromEnvironment(false)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should be possible to fetch the status of control plane pods", func() {
		By("Fetching the cluster status")
		status, err := cli.OpenShiftManagedClustersAdmin.ClusterStatus(context.Background(), os.Getenv("RESOURCEGROUP"), os.Getenv("RESOURCEGROUP"))
		Expect(err).NotTo(HaveOccurred())
		Expect(status).NotTo(BeNil())
	})
})
