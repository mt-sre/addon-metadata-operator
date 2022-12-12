package testutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateRootDirContainsGoMod(t *testing.T) {
	files, err := os.ReadDir(string(RootDir()))
	require.NoError(t, err)

	foundGoMod := false
	for _, file := range files {
		if file.Name() == "go.mod" {
			foundGoMod = true
			break
		}
	}
	require.True(t, foundGoMod)
}
