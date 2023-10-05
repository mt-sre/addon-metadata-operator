//go:build !unit
// +build !unit

package mtcli

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("validate subcommand", func() {
	type namespaceTestCase struct {
		MetadataPath       string
		ExcludedNamespaces []string
		ShouldSucceed      bool
	}

	DescribeTable("AM0008 enabled",
		func(tc namespaceTestCase) {
			cmd := exec.Command(
				_binPath, "validate", "--env", "stage",
				"--enabled", "AM0008",
				"--excluded-namespaces", strings.Join(tc.ExcludedNamespaces, ","),
				tc.MetadataPath,
			)
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
		Entry("reference-addon.v0.0.1-invalid",
			namespaceTestCase{
				MetadataPath:  filepath.Join(testutils.RootDir().TestData().MetadataV1().Legacy(), "reference-addon"),
				ShouldSucceed: false,
			},
		),
		Entry("reference-addon.v0.0.1-valid",
			namespaceTestCase{
				MetadataPath:       filepath.Join(testutils.RootDir().TestData().MetadataV1().Legacy(), "reference-addon"),
				ExcludedNamespaces: []string{"reference-addon"},
				ShouldSucceed:      true,
			},
		),
	)
})
