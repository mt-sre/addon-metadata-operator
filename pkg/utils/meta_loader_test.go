package utils

import (
	"path"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/stretchr/testify/require"
)

var (
	ReferenceAddonImageSetDir   = path.Join(testutils.AddonsImagesetDir, "reference-addon")
	ReferenceAddonIndexImageDir = path.Join(testutils.AddonsIndexImageDir, "reference-addon")
)

// Version is ignored in the case of the static indexImage
func TestMetaLoaderIndexImage(t *testing.T) {
	expectedIndexImage := "quay.io/osd-addons/reference-addon-index@sha256:b9e87a598e7fd6afb4bfedb31e4098435c2105cc8ebe33231c341e515ba9054d"
	cases := []struct {
		addonDir string
		env      string
		version  string
	}{
		{
			addonDir: ReferenceAddonIndexImageDir,
			env:      "stage",
			version:  "latest",
		},
		{
			addonDir: ReferenceAddonIndexImageDir,
			env:      "stage",
			version:  "0.0.1",
		},
	}
	for _, tc := range cases {
		tc := tc // pin
		t.Run(path.Base(tc.addonDir), func(t *testing.T) {
			t.Parallel()
			loader := NewMetaLoader(tc.addonDir, tc.env, tc.version)
			meta, err := loader.Load()
			require.NoError(t, err)
			require.Equal(t, *meta.IndexImage, expectedIndexImage)
			require.Nil(t, meta.ImageSetVersion)
		})
	}
}

// The addonImageSetVersion field is overriden if a version is specified
func TestMetaLoaderImageSet(t *testing.T) {
	expectedImageSetVersion := "0.0.1"
	cases := []struct {
		addonDir           string
		env                string
		version            string
		expectedIndexImage string
	}{
		{
			addonDir:           ReferenceAddonImageSetDir,
			env:                "stage",
			version:            "latest", // 0.0.2
			expectedIndexImage: "quay.io/osd-addons/reference-addon-index@sha256:d9f95ecd3cace47d9f34f63354c06d829abd165de90ef3990b6f9feda36233f3",
		},
		{
			addonDir:           ReferenceAddonImageSetDir,
			env:                "stage",
			version:            "0.0.1",
			expectedIndexImage: "quay.io/osd-addons/reference-addon-index@sha256:b9e87a598e7fd6afb4bfedb31e4098435c2105cc8ebe33231c341e515ba9054d",
		},
		{
			addonDir:           ReferenceAddonImageSetDir,
			env:                "stage",
			version:            "", // defaults to addonImageSetVersion == 0.0.1
			expectedIndexImage: "quay.io/osd-addons/reference-addon-index@sha256:b9e87a598e7fd6afb4bfedb31e4098435c2105cc8ebe33231c341e515ba9054d",
		},
	}
	for _, tc := range cases {
		tc := tc // pin
		t.Run(path.Base(tc.addonDir), func(t *testing.T) {
			t.Parallel()
			loader := NewMetaLoader(tc.addonDir, tc.env, tc.version)
			meta, err := loader.Load()
			require.NoError(t, err)
			require.Equal(t, *meta.IndexImage, tc.expectedIndexImage)
			// the expectedImageSetVersion is always 0.0.1, because this is the
			// value set in the manifest. The metaLoader can override this
			// default value, for example when setting latest.
			require.Equal(t, *meta.ImageSetVersion, expectedImageSetVersion)
		})
	}
}
