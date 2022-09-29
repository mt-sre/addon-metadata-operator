package extractor

import (
	"context"

	"github.com/mt-sre/addon-metadata-operator/pkg/operator"
)

// Extractor - utilizes both the indexExtractor and bundleExtractor to first extract
// all bundleImages from an indexImage, and then extract all the corresponding
// bundles from those underlying bundleImages.
type Extractor interface {
	// extract all bundles from indexImage matching pkgName
	ExtractBundles(indexImage string, pkgName string) ([]operator.Bundle, error)
	// extract all bundles from indexImage, for all packages
	ExtractAllBundles(indexImage string) ([]operator.Bundle, error)
}

// Choosing to list explicitly that the indexExtractor works with BundleImages
// as it is not obvious from it's underlying name.

// IndexExtractor - extracts bundleImages from an indexImage. Supports both the sql
// and file based catalog format. An indexImage contains one or multiple packages,
// which contain bundleImages.
// Catalog format: https://docs.openshift.com/container-platform/4.9/operators/admin/olm-managing-custom-catalogs.html#olm-managing-custom-catalogs-fb
type IndexExtractor interface {
	// extract all bundleImages contained in indexImage, matching pkgName
	ExtractBundleImages(indexImage string, pkgName string) ([]string, error)
	// extract all bundleImages contained in the indexImage, for all packages.
	ExtractAllBundleImages(indexImage string) ([]string, error)
}

// BundleExtractor - extracts a single bundle from it's bundleImage, using the bundle
// format by OPM.
// Bundle format: https://docs.openshift.com/container-platform/4.9/operators/understanding/olm-packaging-format.html#olm-bundle-format_olm-packaging-format
type BundleExtractor interface {
	Extract(ctx context.Context, bundleImage string) (operator.Bundle, error)
}
