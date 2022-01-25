package extractor

import (
	"testing"

	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/stretchr/testify/require"
)

func TestMainExtractorImplements(t *testing.T) {
	require.Implements(t, new(Extractor), &MainExtractor{})
}

func TestMainExtractorWithDefaultValues(t *testing.T) {
	cases := []struct {
		indexImage           string
		pkgName              string
		expectedBundleImages map[string]bool
	}{
		{
			// sql-based catalog image
			indexImage: "quay.io/osd-addons/reference-addon-index@sha256:b9e87a598e7fd6afb4bfedb31e4098435c2105cc8ebe33231c341e515ba9054d",
			pkgName:    "",
			expectedBundleImages: map[string]bool{
				"quay.io/osd-addons/reference-addon-bundle:0.1.0-c15cedb": true,
				"quay.io/osd-addons/reference-addon-bundle:0.1.1-c15cedb": true,
				"quay.io/osd-addons/reference-addon-bundle:0.1.2-c15cedb": true,
				"quay.io/osd-addons/reference-addon-bundle:0.1.3-c15cedb": true,
				"quay.io/osd-addons/reference-addon-bundle:0.1.4-c15cedb": true,
				"quay.io/osd-addons/reference-addon-bundle:0.1.5-c15cedb": true,
			},
		},
		{
			// file-based catalog image
			indexImage: "quay.io/sblaisdo/reference-addon-index:test",
			pkgName:    "reference-addon",
			expectedBundleImages: map[string]bool{
				"quay.io/osd-addons/reference-addon-bundle:0.1.0-bcb6192": true,
				"quay.io/osd-addons/reference-addon-bundle:0.1.1-bcb6192": true,
				"quay.io/osd-addons/reference-addon-bundle:0.1.2-bcb6192": true,
				"quay.io/osd-addons/reference-addon-bundle:0.1.3-bcb6192": true,
				"quay.io/osd-addons/reference-addon-bundle:0.1.4-bcb6192": true,
				"quay.io/osd-addons/reference-addon-bundle:0.1.5-bcb6192": true,
				"quay.io/osd-addons/reference-addon-bundle:0.1.6-bcb6192": true,
			},
		},
	}
	extractor := New()
	for _, tc := range cases {
		tc := tc // pin
		t.Run(tc.indexImage, func(t *testing.T) {
			t.Parallel()

			var bundles []*registry.Bundle
			var err error

			if tc.pkgName == "" {
				bundles, err = extractor.ExtractAllBundles(tc.indexImage)
			} else {
				bundles, err = extractor.ExtractBundles(tc.indexImage, tc.pkgName)
			}
			require.NoError(t, err)

			require.Equal(t, len(bundles), len(tc.expectedBundleImages))
			for _, bundle := range bundles {
				_, ok := tc.expectedBundleImages[bundle.BundleImage]
				require.True(t, ok)
			}
		})
	}
}
