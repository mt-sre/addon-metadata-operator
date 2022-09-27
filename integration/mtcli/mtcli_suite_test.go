package mtcli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	_binPath string
)

func TestMTCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mtcli suite")
}

var _ = BeforeSuite(func() {
	root, err := projectRoot()
	Expect(err).ToNot(HaveOccurred())

	env := []string{
		"CGO_ENABLED=1",
		"CGO_CFLAGS=-DSQLITE_ENABLE_JSON1",
	}

	_binPath, err = gexec.BuildWithEnvironment(filepath.Join(root, "cmd", "mtcli"), env)
	Expect(err).ToNot(HaveOccurred())

	DeferCleanup(gexec.CleanupBuildArtifacts)
})

var errSetup = errors.New("test setup failed")

func projectRoot() (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Stdout = &buf
	cmd.Stderr = io.Discard

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("determining top level directory from git: %w", errSetup)
	}

	return strings.TrimSpace(buf.String()), nil
}
