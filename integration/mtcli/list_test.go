//go:build !unit
// +build !unit

package mtcli

import (
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("list subcommand", func() {
	type bundlesTestCase struct {
		IndexImage      string
		ExpectedBundles []string
	}

	DescribeTable("bundles subcommand",
		func(tc bundlesTestCase) {
			cmd := exec.Command(_binPath, "list", "bundles", tc.IndexImage)

			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, "5s").Should(Exit(0))

			Expect(session.Out).To(Say(strings.Join(tc.ExpectedBundles, "\n")))
		},
		Entry("reference-addon v0.1.5",
			bundlesTestCase{
				IndexImage: "quay.io/osd-addons/reference-addon-index@sha256:b9e87a598e7fd6afb4bfedb31e4098435c2105cc8ebe33231c341e515ba9054d",
				ExpectedBundles: []string{
					"reference-addon.v0.1.0",
					"reference-addon.v0.1.1",
					"reference-addon.v0.1.2",
					"reference-addon.v0.1.3",
					"reference-addon.v0.1.4",
					"reference-addon.v0.1.5",
				},
			},
		),
		Entry("ocs-converged v0.1.1",
			bundlesTestCase{
				IndexImage: "quay.io/osd-addons/ocs-converged-index@sha256:24c6519b0d109a8e1e5349706a95d05e268a74f7df8f9040fc3a805700169afe",
				ExpectedBundles: []string{
					"ocs-osd-deployer.v1.0.0",
					"ocs-osd-deployer.v1.0.1",
					"ocs-osd-deployer.v1.0.2",
					"ocs-osd-deployer.v1.1.0",
					"ocs-osd-deployer.v1.1.1",
					"ose-prometheus-operator.4.8.0",
				},
			},
		),
	)
})
