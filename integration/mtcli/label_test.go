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
	type labelTestCase struct {
		MetadataPath  string
		ShouldSucceed bool
	}

	DescribeTable("AM0002 enabled",
		func(tc labelTestCase) {
			cmd := exec.Command(_binPath, "validate", "--env", "stage", "--enabled", "AM0002", tc.MetadataPath)
			cmd.Env = []string{
				`OCM_TOKEN=""`,
			}

			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			exitCode := 0
			if !tc.ShouldSucceed {
				exitCode = 1
			}

			Eventually(session, "15s").Should(Exit(exitCode))
		},
		Entry("reference-addon.v0.0.1-valid",
			labelTestCase{
				MetadataPath:  filepath.Join(testutils.RootDir().TestData().MetadataV1().Legacy(), "reference-addon"),
				ShouldSucceed: true,
			},
		),
		Entry("ocm-addon-test-operator.v0.0.12-valid",
			labelTestCase{
				MetadataPath:  filepath.Join(testutils.RootDir().TestData().MetadataV1().Legacy(), "ocm-addon-test-operator"),
				ShouldSucceed: true,
			},
		),
		Entry("connectors-operator.v1.1.17-invalid",
			labelTestCase{
				MetadataPath:  filepath.Join(testutils.RootDir().TestData().MetadataV1().Legacy(), "connectors-operator"),
				ShouldSucceed: false,
			},
		),
	)
})
