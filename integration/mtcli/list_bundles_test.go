//go:build !unit
// +build !unit

package mtcli_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type e2eTestSuite struct {
	suite.Suite
	mtcliPath string
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) TearDownSuite() {
	// Remove cache directory
	err := os.RemoveAll(utils.DefaultCacheDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to cleanup cache dirs")
	}
}

func (s *e2eTestSuite) SetupSuite() {
	s.mtcliPath = os.Getenv("E2E_MTCLI_PATH")
	_, err := os.Stat(s.mtcliPath)
	s.Require().NoError(err, "Cant find mtcli binary at E2E_MTCLI_PATH")
}

func (s *e2eTestSuite) TestMtcliListBundlesE2E() {
	testCases := []struct {
		indexImageUrl  string
		expectedOutput []string
	}{
		{
			indexImageUrl: "quay.io/osd-addons/reference-addon-index@sha256:b9e87a598e7fd6afb4bfedb31e4098435c2105cc8ebe33231c341e515ba9054d",
			expectedOutput: []string{
				"reference-addon.v0.1.0",
				"reference-addon.v0.1.1",
				"reference-addon.v0.1.2",
				"reference-addon.v0.1.3",
				"reference-addon.v0.1.4",
				"reference-addon.v0.1.5",
			},
		},
		{
			indexImageUrl: "quay.io/osd-addons/ocs-converged-index@sha256:24c6519b0d109a8e1e5349706a95d05e268a74f7df8f9040fc3a805700169afe",
			expectedOutput: []string{
				"ocs-osd-deployer.v1.0.0",
				"ocs-osd-deployer.v1.0.1",
				"ocs-osd-deployer.v1.0.2",
				"ocs-osd-deployer.v1.1.0",
				"ocs-osd-deployer.v1.1.1",
				"ose-prometheus-operator.4.8.0",
			},
		},
	}
	for _, testCase := range testCases {
		cmd := prepareListBundlesCmd(s.mtcliPath, testCase.indexImageUrl)
		outBytes, err := cmd.CombinedOutput()
		s.Require().NoError(err)
		// remove last trailing newline
		outString := strings.TrimSuffix(string(outBytes), "\n")
		outLines := strings.Split(outString, "\n")
		s.Equal(testCase.expectedOutput, outLines)

		// ensure cache keys
		cacheContents, err := readCacheFile(testCase.indexImageUrl)
		s.Require().NoError(err)
		s.Equal(cacheContents, expectedCacheValue(testCase.indexImageUrl))
	}
}

func expectedCacheValue(indexImageUrl string) string {
	return strings.Join(
		[]string{indexImageUrl, utils.AllAddonsIdentifier},
		"<>",
	)
}

func readCacheFile(indexImageUrl string) (string, error) {
	fileBytes, err := ioutil.ReadFile(filepath.Join(
		utils.DefaultCacheDir,
		indexImageUrl,
		utils.DefaultCacheFileName,
	))
	if err != nil {
		return "", err
	}
	return string(fileBytes), nil
}

func prepareListBundlesCmd(cliPath string, arg string) *exec.Cmd {
	return exec.Command(cliPath, "list", "bundles", arg)
}
