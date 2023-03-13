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

var _ = Describe("validate subcommand", func() {
	type disableTestCase struct {
		MetadataPath  string
		ShouldSucceed bool
	}

	DescribeTable("AM0002, AM0004, AM0011, AM0015 disabled",
		func(tc disableTestCase) {
			cmd := exec.Command(_binPath, "validate", "--env", "stage", "--disabled", "AM0002,AM0004,AM0011,AM0015", tc.MetadataPath)
			cmd.Env = []string{
				`OCM_TOKEN=""`,
			}

			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			exitCode := 0
			if !tc.ShouldSucceed {
				exitCode = 1
			}

			Eventually(session, "30s").Should(Exit(exitCode))
		},
		Entry("ocm-addon-test-operator.v0.0.12-invalid",
			disableTestCase{
				MetadataPath:  filepath.Join(testutils.RootDir().TestData().MetadataV1().Legacy(), "ocm-addon-test-operator"),
				ShouldSucceed: false,
			},
		),
		Entry("reference-addon.v0.0.5-imagesets-valid",
			disableTestCase{
				MetadataPath:  filepath.Join(testutils.RootDir().TestData().MetadataV1().ImageSets(), "reference-addon"),
				ShouldSucceed: true,
			},
		),
		Entry("advanced-cluster-management.v2.7.0-invalid",
			disableTestCase{
				MetadataPath:  filepath.Join(testutils.RootDir().TestData().MetadataV1().Legacy(), "advanced-cluster-management"),
				ShouldSucceed: false,
			},
		),
	)
})
