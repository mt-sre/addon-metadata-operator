package extractor

import (
	"context"

	"github.com/operator-framework/operator-registry/pkg/registry"
)

// Extractor - utilizes both the indexExtractor and bundleExtractor to first extract
// all bundleImages from an indexImage, and then extract all the corresponding
// bundles from those underlying bundleImages.
type Extractor interface {
	// extract all bundles from indexImage matching pkgName
	ExtractBundles(indexImage string, pkgName string) ([]*registry.Bundle, error)
	// extract all bundles from indexImage, for all packages
	ExtractAllBundles(indexImage string) ([]*registry.Bundle, error)
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

// Cache - caches bundleImages using a double index. We first cache results per
// indexImage and then by pkgName.
type IndexCache interface {
	GetBundleImages(indexImage string, pkgName string) []string
	// An indexImage is mutable, so we have to update the cache on every call.
	SetBundleImages(indexImage string, bundleImagesMap map[string][]string)
}

// BundleExtractor - extracts a single bundle from it's bundleImage, using the bundle
// format by OPM.
// Bundle format: https://docs.openshift.com/container-platform/4.9/operators/understanding/olm-packaging-format.html#olm-bundle-format_olm-packaging-format
type BundleExtractor interface {
	Extract(ctx context.Context, bundleImage string) (*registry.Bundle, error)
}

// BundleCache - cache a bundle object to avoid having to pull and extract the bundleImage.
type BundleCache interface {
	Get(bundleImage string) (*registry.Bundle, error)
	// bundles are defined to be immutable by OPM, se we don't have to update
	// the cache if an entry already exists.
	Set(bundleImage string, bundle *registry.Bundle) error
}

// BundleEncoder - encodes a bundle object for more efficient caching. Bundle objects
// are large structure and we want to compress and work with a friendlier format
// when caching or sending them on the wire (e.g.: Redis)
type BundleEncoder interface {
	Encode(bundle *registry.Bundle) ([]byte, error)
	Decode(data []byte) (*registry.Bundle, error)
}
