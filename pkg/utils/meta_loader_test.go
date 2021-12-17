package utils_test

import (
	"fmt"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/stretchr/testify/require"
)

// Version is ignored in the case of the static indexImage
func TestMetaLoaderStaticIndexImage(t *testing.T) {
	env := "stage"
	refAddonStage, err := testutils.GetReferenceAddonStage()
	require.NoError(t, err)

	versions := []string{"latest", "0.0.1"}
	for _, version := range versions {
		version := version // pin
		t.Run(fmt.Sprintf("load-reference-addon-indeximage-version-%s", version), func(t *testing.T) {
			t.Parallel()
			loader := utils.NewMetaLoader(refAddonStage.IndexImageDir(), env, version)
			meta, err := loader.Load()
			require.NoError(t, err)
			require.Equal(t, *meta.IndexImage, *refAddonStage.MetaIndexImage.IndexImage)
			require.Nil(t, meta.ImageSetVersion)
		})
	}
}

// The addonImageSetVersion field is overriden if a different version is specified
// Three imageset cases
// 1. "latest" resolves to latest imageset in reference-addon/addonimagesets/stage/*.yaml
// 2. "0.0.1"  resolves to reference-addon/addonimagesets/stage/reference-addon.v0.0.1.yaml
// 3. ""       uses value from reference-addon/metadata/stage/addon.yaml::addonImageSetVersion
func TestMetaLoaderImageSet(t *testing.T) {
	env := "stage"
	refAddonStage, err := testutils.GetReferenceAddonStage()
	require.NoError(t, err)

	cases := []struct {
		version         string
		expectedVersion string
	}{
		{
			version:         "latest",
			expectedVersion: "0.0.5",
		},
		{
			version:         "0.0.1",
			expectedVersion: "0.0.1",
		},
		{
			version:         "",
			expectedVersion: "0.0.5",
		},
	}
	for _, tc := range cases {
		tc := tc // pin
		t.Run(fmt.Sprintf("load-reference-addon-imageset-version-%s", tc), func(t *testing.T) {
			t.Parallel()
			loader := utils.NewMetaLoader(refAddonStage.ImageSetDir(), env, tc.version)
			meta, err := loader.Load()
			require.NoError(t, err)

			expectedImageSet, err := refAddonStage.GetImageSet(tc.version)
			require.NoError(t, err)

			expectedImageSetVersion, err := expectedImageSet.GetSemver()
			require.NoError(t, err)

			require.Equal(t, *meta.IndexImage, expectedImageSet.IndexImage)
			require.Equal(t, *meta.ImageSetVersion, expectedImageSetVersion)

			// Also do a static check, to bullet proof code from testutils as well
			// If any failures happen here, prompt the user to update the manifests
			errorMsg := `
			Did you update manifests in internal/testdata/addons-imageset/reference-addon/ ?
			Please make sure you update the static "expectedVersion" fields in this test to match your changes.
			`
			require.Equal(t, tc.expectedVersion, expectedImageSetVersion, errorMsg)
			require.Equal(t, tc.expectedVersion, *meta.ImageSetVersion, errorMsg)
		})
	}
}
