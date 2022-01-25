package extractor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultIndexExtractorImplements(t *testing.T) {
	require.Implements(t, new(IndexExtractor), &DefaultIndexExtractor{})
}

func TestExtractorFileBasedAndSQLCatalogs(t *testing.T) {
	cases := []struct {
		indexImage           string
		pkgName              string
		expectedBundleImages []string
	}{
		{
			// sql-based catalog image
			indexImage: "quay.io/osd-addons/reference-addon-index@sha256:b9e87a598e7fd6afb4bfedb31e4098435c2105cc8ebe33231c341e515ba9054d",
			pkgName:    "reference-addon",
			expectedBundleImages: []string{
				"quay.io/osd-addons/reference-addon-bundle:0.1.0-c15cedb",
				"quay.io/osd-addons/reference-addon-bundle:0.1.1-c15cedb",
				"quay.io/osd-addons/reference-addon-bundle:0.1.2-c15cedb",
				"quay.io/osd-addons/reference-addon-bundle:0.1.3-c15cedb",
				"quay.io/osd-addons/reference-addon-bundle:0.1.4-c15cedb",
				"quay.io/osd-addons/reference-addon-bundle:0.1.5-c15cedb",
			},
		},
		{
			// file-based catalog image
			indexImage: "quay.io/sblaisdo/reference-addon-index:test",
			pkgName:    "reference-addon",
			expectedBundleImages: []string{
				"quay.io/osd-addons/reference-addon-bundle:0.1.0-bcb6192",
				"quay.io/osd-addons/reference-addon-bundle:0.1.1-bcb6192",
				"quay.io/osd-addons/reference-addon-bundle:0.1.2-bcb6192",
				"quay.io/osd-addons/reference-addon-bundle:0.1.3-bcb6192",
				"quay.io/osd-addons/reference-addon-bundle:0.1.4-bcb6192",
				"quay.io/osd-addons/reference-addon-bundle:0.1.5-bcb6192",
				"quay.io/osd-addons/reference-addon-bundle:0.1.6-bcb6192",
			},
		},
	}
	cache := NewIndexMemoryCache()
	extractor := NewIndexExtractor(WithIndexCache(cache))
	for _, tc := range cases {
		tc := tc // pin
		t.Run(tc.indexImage, func(t *testing.T) {
			t.Parallel()
			bundleImages, err := extractor.ExtractBundleImages(tc.indexImage, tc.pkgName)
			require.NoError(t, err)
			require.ElementsMatch(t, bundleImages, tc.expectedBundleImages)

			cachedBundleImages := cache.GetBundleImages(tc.indexImage, tc.pkgName)
			require.ElementsMatch(t, cachedBundleImages, tc.expectedBundleImages)
		})
	}
}

func TestIndexExtractorListAllBundles(t *testing.T) {
	t.Parallel()
	indexImage := "quay.io/osd-addons/gpu-operator-index@sha256:62e0f330276375758f875c62c90e6c3e4e217247f221c96ce5e4ab64f6617e38"
	bundleImagesMap := map[string][]string{
		"gpu-operator-certified-addon": {
			"quay.io/osd-addons/gpu-operator-bundle:1.7.1-0ddc381",
			"quay.io/osd-addons/gpu-operator-bundle:1.8.0-0ddc381",
			"quay.io/osd-addons/gpu-operator-bundle:1.8.2-0ddc381",
			"quay.io/osd-addons/gpu-operator-bundle:1.8.3-0ddc381",
			"quay.io/osd-addons/gpu-operator-bundle:1.9.0-beta-0ddc381",
		},
		"node-feature-discovery-operator": {
			"quay.io/osd-addons/gpu-operator-nfd-operator-bundle:4.8.0-0ddc381",
		},
	}
	expectedBundleImages := []string{
		"quay.io/osd-addons/gpu-operator-bundle:1.7.1-0ddc381",
		"quay.io/osd-addons/gpu-operator-bundle:1.8.0-0ddc381",
		"quay.io/osd-addons/gpu-operator-bundle:1.8.2-0ddc381",
		"quay.io/osd-addons/gpu-operator-bundle:1.8.3-0ddc381",
		"quay.io/osd-addons/gpu-operator-bundle:1.9.0-beta-0ddc381",
		"quay.io/osd-addons/gpu-operator-nfd-operator-bundle:4.8.0-0ddc381",
	}
	cache := NewIndexMemoryCache()
	extractor := NewIndexExtractor(WithIndexCache(cache))

	// test bundleImages for all packages are listed
	bundleImages, err := extractor.ExtractAllBundleImages(indexImage)
	require.NoError(t, err)
	require.ElementsMatch(t, bundleImages, expectedBundleImages)

	cachedBundleImages := cache.GetBundleImages(indexImage, allBundlesKey)
	require.ElementsMatch(t, cachedBundleImages, expectedBundleImages)

	// test bundles have been cached per pkgName
	for pkgName, expectedBundleImages := range bundleImagesMap {
		bundleImages, err := extractor.ExtractBundleImages(indexImage, pkgName)
		require.NoError(t, err)
		require.ElementsMatch(t, bundleImages, expectedBundleImages)

		cachedBundleImages := cache.GetBundleImages(indexImage, pkgName)
		require.ElementsMatch(t, cachedBundleImages, expectedBundleImages)
	}
}
