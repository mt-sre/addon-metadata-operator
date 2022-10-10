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
	type testHarnessTestCase struct {
		MetadataPath  string
		ShouldSucceed bool
	}

	DescribeTable("AM0005 enabled",
		func(tc testHarnessTestCase) {
			cmd := exec.Command(_binPath, "validate", "--env", "stage", "--enabled", "AM0005", tc.MetadataPath)
			cmd.Env = []string{
				`OCM_TOKEN=""`,
			}

			session, err := Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			exitCode := 0
			if !tc.ShouldSucceed {
				exitCode = 1
			}

			Eventually(session, "10s").Should(Exit(exitCode))
		},
		Entry("reference-addon.v0.0.1-valid",
			testHarnessTestCase{
				MetadataPath:  filepath.Join(testutils.RootDir().TestData().MetadataV1().Legacy(), "reference-addon"),
				ShouldSucceed: true,
			},
		),
		Entry("connectors-operator.v1.1.17-invalid",
			testHarnessTestCase{
				MetadataPath:  filepath.Join(testutils.RootDir().TestData().MetadataV1().Legacy(), "connectors-operator"),
				ShouldSucceed: false,
			},
		),
		Entry("ocm-addon-test-operator.v0.0.12-valid",
			testHarnessTestCase{
				MetadataPath:  filepath.Join(testutils.RootDir().TestData().MetadataV1().Legacy(), "ocm-addon-test-operator"),
				ShouldSucceed: true,
			},
		),
		Entry("reference-addon.v0.0.5-imagesets-valid",
			testHarnessTestCase{
				MetadataPath:  filepath.Join(testutils.RootDir().TestData().MetadataV1().ImageSets(), "reference-addon"),
				ShouldSucceed: true,
			},
		),
	)
})
