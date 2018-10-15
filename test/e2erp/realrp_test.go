//+build e2erp

package e2erp

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resource provider e2e tests [Real]", func() {
	defer GinkgoRecover()

	It("dummy real test", func() {
		var err error
		Expect(err).NotTo(HaveOccurred())
	})

	It("should not be possible for customer to mutate an osa scale set", func() {
		var err error
		Expect(err).NotTo(HaveOccurred())

		By("Updating the instance count")

		Expect(err).NotTo(HaveOccurred())

		By("Updating the instance type")

		Expect(err).NotTo(HaveOccurred())

		By("Updating the ssh key")

		Expect(err).NotTo(HaveOccurred())

		By("Rebooting instances")

		Expect(err).NotTo(HaveOccurred())

		By("Creating script extensions")

		Expect(err).NotTo(HaveOccurred())

	})
})
