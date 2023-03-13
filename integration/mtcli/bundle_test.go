//go:build !unit
// +build !unit

package mtcli

import (
	"os/exec"
	"path/filepath"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("bundle subcommand", func() {
	type validateTestCase struct {
		BundlePath    string
		ShouldSucceed bool
	}

	DescribeTable("validate subcommand",
		func(tc validateTestCase) {
			cmd := exec.Command(_binPath, "bundle", "validate", tc.BundlePath)

			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			exitCode := 0
			if !tc.ShouldSucceed {
				exitCode = 1
			}

			Eventually(session, "30s").Should(Exit(exitCode))
		},
		Entry("reference-addon.0.1.6-valid",
			validateTestCase{
				BundlePath:    filepath.Join(testutils.RootDir().TestData().Bundles(), "reference-addon", "main", "0.1.6"),
				ShouldSucceed: true,
			},
		),
		Entry("addon-operator.0.3.0-valid",
			validateTestCase{
				BundlePath:    filepath.Join(testutils.RootDir().TestData().Bundles(), "reference-addon", "addon-operator", "0.3.0"),
				ShouldSucceed: true,
			},
		),
		Entry("rhods.1.1.1-57-invalid",
			validateTestCase{
				BundlePath:    filepath.Join(testutils.RootDir().TestData().Bundles(), "rhods", "main", "1.1.1-57"),
				ShouldSucceed: false,
			},
		),
	)
})
