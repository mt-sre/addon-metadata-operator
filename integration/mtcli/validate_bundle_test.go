//go:build !unit
// +build !unit

package mtcli_test

import (
	"os/exec"
	"path/filepath"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
)

func (s *e2eTestSuite) TestMtcliValidateBundle() {
	cases := []struct {
		name       string
		bundlePath string
		isValid    bool
	}{
		{
			name:       "reference-addon.0.1.6-valid",
			bundlePath: filepath.Join(testutils.RootDir().TestData().Bundles(), "reference-addon", "main", "0.1.6"),
			isValid:    true,
		},
		{
			name:       "addon-operator.0.3.0-valid",
			bundlePath: filepath.Join(testutils.RootDir().TestData().Bundles(), "reference-addon", "addon-operator", "0.3.0"),
			isValid:    true,
		},
		{
			name:       "rhods.1.1.1-57-invalid",
			bundlePath: filepath.Join(testutils.RootDir().TestData().Bundles(), "rhods", "main", "1.1.1-57"),
			isValid:    false,
		},
	}
	for _, tc := range cases {
		cmd := exec.Command(s.mtcliPath, "bundle", "validate", tc.bundlePath)
		err := cmd.Run() // we don't care about the output

		if tc.isValid {
			s.Require().NoError(err)
		} else {
			s.Require().NotNil(err)
		}
	}
}
