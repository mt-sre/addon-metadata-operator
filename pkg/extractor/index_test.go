package extractor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndexExtractorFileBasedAndSQLCatalogs(t *testing.T) {
	cases := []struct {
		indexImage      string
		pkgName         string
		expectedBundles []string
	}{
		{
			// sql-based catalog image
			indexImage: "quay.io/osd-addons/reference-addon-index@sha256:b9e87a598e7fd6afb4bfedb31e4098435c2105cc8ebe33231c341e515ba9054d",
			pkgName:    "reference-addon",
			expectedBundles: []string{
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
			expectedBundles: []string{
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
	extractor := NewIndexExtractor(cache)
	for _, tc := range cases {
		tc := tc // pin
		t.Run(tc.indexImage, func(t *testing.T) {
			t.Parallel()
			bundles, err := extractor.ListBundlesFromPackage(tc.indexImage, tc.pkgName)
			require.NoError(t, err)
			require.Equal(t, bundles, tc.expectedBundles)
			cachedBundles := cache.GetBundles(tc.indexImage, tc.pkgName)
			require.Equal(t, cachedBundles, tc.expectedBundles)
		})
	}
}

// quay.io/osd-addons/gpu-operator-index@sha256:62e0f330276375758f875c62c90e6c3e4e217247f221c96ce5e4ab64f6617e38
func TestIndexExtractorListAllBundles(t *testing.T) {
	t.Parallel()
	indexImage := "quay.io/osd-addons/gpu-operator-index@sha256:62e0f330276375758f875c62c90e6c3e4e217247f221c96ce5e4ab64f6617e38"
	pkgBundlesMap := map[string][]string{
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
	expectedAllBundles := []string{
		"quay.io/osd-addons/gpu-operator-bundle:1.7.1-0ddc381",
		"quay.io/osd-addons/gpu-operator-bundle:1.8.0-0ddc381",
		"quay.io/osd-addons/gpu-operator-bundle:1.8.2-0ddc381",
		"quay.io/osd-addons/gpu-operator-bundle:1.8.3-0ddc381",
		"quay.io/osd-addons/gpu-operator-bundle:1.9.0-beta-0ddc381",
		"quay.io/osd-addons/gpu-operator-nfd-operator-bundle:4.8.0-0ddc381",
	}
	cache := NewIndexMemoryCache()
	extractor := NewIndexExtractor(cache)

	// test bundles for all packages are listed
	bundles, err := extractor.ListAllBundles(indexImage)
	require.NoError(t, err)
	require.Equal(t, bundles, expectedAllBundles)
	allCachedBundles := cache.GetBundles(indexImage, allBundlesKey)
	require.Equal(t, allCachedBundles, expectedAllBundles)

	// test bundles have been cached per pkgName
	for pkgName, expectedBundles := range pkgBundlesMap {
		bundles, err := extractor.ListBundlesFromPackage(indexImage, pkgName)
		require.NoError(t, err)
		require.Equal(t, bundles, expectedBundles)
		cachedBundles := cache.GetBundles(indexImage, pkgName)
		require.Equal(t, cachedBundles, expectedBundles)
	}
}
