package extractor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultIndexExtractorImplements(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(IndexExtractor), &DefaultIndexExtractor{})
}

func TestExtractorFileBasedAndSQLCatalogs(t *testing.T) {
	t.Parallel()

	cache := NewIndexCacheImpl()
	extractor := NewIndexExtractor(WithIndexCache(cache))

	for name, tc := range map[string]struct {
		IndexImage           string
		PkgName              string
		ExpectedBundleImages []string
	}{
		"sql-based catalog image": {
			IndexImage: "quay.io/osd-addons/reference-addon-index@sha256:b9e87a598e7fd6afb4bfedb31e4098435c2105cc8ebe33231c341e515ba9054d",
			PkgName:    "reference-addon",
			ExpectedBundleImages: []string{
				"quay.io/osd-addons/reference-addon-bundle:0.1.0-c15cedb",
				"quay.io/osd-addons/reference-addon-bundle:0.1.1-c15cedb",
				"quay.io/osd-addons/reference-addon-bundle:0.1.2-c15cedb",
				"quay.io/osd-addons/reference-addon-bundle:0.1.3-c15cedb",
				"quay.io/osd-addons/reference-addon-bundle:0.1.4-c15cedb",
				"quay.io/osd-addons/reference-addon-bundle:0.1.5-c15cedb",
			},
		},
		"file-based catalog image": {
			IndexImage: "quay.io/osd-addons/reference-addon-index:file-based-poc",
			PkgName:    "reference-addon",
			ExpectedBundleImages: []string{
				"quay.io/osd-addons/reference-addon-bundle:0.1.6-single",
			},
		},
	} {
		tc := tc // pin

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bundleImages, err := extractor.ExtractBundleImages(context.Background(), tc.IndexImage, tc.PkgName)
			require.NoError(t, err)

			assert.ElementsMatch(t, bundleImages, tc.ExpectedBundleImages)

			cachedBundleImages, err := cache.GetBundleImages(tc.IndexImage, tc.PkgName)
			require.NoError(t, err)

			assert.ElementsMatch(t, cachedBundleImages, tc.ExpectedBundleImages)
		})
	}
}

func TestIndexExtractorListAllBundles(t *testing.T) {
	t.Parallel()

	cache := NewIndexCacheImpl()
	extractor := NewIndexExtractor(WithIndexCache(cache))

	for name, tc := range map[string]struct {
		IndexImage           string
		PackageToBundle      map[string][]string
		ExpectedBundleImages []string
	}{
		"gpu-operator": {
			IndexImage: "quay.io/osd-addons/gpu-operator-index@sha256:62e0f330276375758f875c62c90e6c3e4e217247f221c96ce5e4ab64f6617e38",
			PackageToBundle: map[string][]string{
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
			},
			ExpectedBundleImages: []string{
				"quay.io/osd-addons/gpu-operator-bundle:1.7.1-0ddc381",
				"quay.io/osd-addons/gpu-operator-bundle:1.8.0-0ddc381",
				"quay.io/osd-addons/gpu-operator-bundle:1.8.2-0ddc381",
				"quay.io/osd-addons/gpu-operator-bundle:1.8.3-0ddc381",
				"quay.io/osd-addons/gpu-operator-bundle:1.9.0-beta-0ddc381",
				"quay.io/osd-addons/gpu-operator-nfd-operator-bundle:4.8.0-0ddc381",
			},
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			// test bundleImages for all packages are listed
			bundleImages, err := extractor.ExtractAllBundleImages(ctx, tc.IndexImage)
			require.NoError(t, err)

			assert.ElementsMatch(t, bundleImages, tc.ExpectedBundleImages)

			cachedBundleImages, err := cache.GetBundleImages(tc.IndexImage, allBundlesKey)
			require.NoError(t, err)

			assert.ElementsMatch(t, cachedBundleImages, tc.ExpectedBundleImages)

			// test bundles have been cached per pkgName
			for pkgName, expectedBundleImages := range tc.PackageToBundle {
				bundleImages, err := extractor.ExtractBundleImages(ctx, tc.IndexImage, pkgName)
				require.NoError(t, err)

				assert.ElementsMatch(t, bundleImages, expectedBundleImages)

				cachedBundleImages, err := cache.GetBundleImages(tc.IndexImage, pkgName)
				require.NoError(t, err)

				assert.ElementsMatch(t, cachedBundleImages, expectedBundleImages)
			}
		})
	}
}
