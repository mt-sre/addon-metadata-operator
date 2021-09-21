package v1alpha1

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestAddonMetadataFromYAML(t *testing.T) {
	files := []string{path.Join(testutils.TestdataDir, "reference-addon.yaml")}
	for _, f := range files {
		f := f
		t.Run(path.Base(f), func(t *testing.T) {
			t.Parallel()
			data, err := ioutil.ReadFile(f)
			require.Nil(t, err)
			addonMetadata := &AddonMetadata{}
			require.Nil(t, addonMetadata.FromYAML(data))
		})
	}
}
