package extractor_test

import (
	"path"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/stretchr/testify/require"
)

/*
	This is to experiment what is most efficient between []*registry.Bundle and
	[]registry.Bundle. As the registry.Bundle struct is massive, it is important
	to benchmark our various use cases.
*/

const (
	benchmarkBundleName = "benchmarkBundle"
	nBundles            = 10
)

var bundlePath = path.Join(testutils.TestdataDir(), "assets/am0007/csv.yaml")

type metaBundleByReference struct {
	Bundles []*registry.Bundle
}

// benchmark bundles workflow using *registry.Bundle{}
func BenchmarkBundlesByReference(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		bundles := make([]*registry.Bundle, nBundles) // pre-allocate
		for j := 0; j < nBundles; j++ {
			bundle, err := testutils.NewBundle(benchmarkBundleName, bundlePath)
			require.NoError(b, err)
			bundles[j] = bundle
		}
		require.Equal(b, len(bundles), nBundles)
		mb := &metaBundleByReference{Bundles: bundles}
		handleBundlesByReference(b, mb)
	}
}

func handleBundlesByReference(b *testing.B, mb *metaBundleByReference) {
	for _, bundle := range mb.Bundles {
		bundleName := getBundleNameByReference(bundle)
		require.Equal(b, bundleName, benchmarkBundleName)
	}
}

func getBundleNameByReference(bundle *registry.Bundle) string {
	return bundle.Name
}

type metaBundleByValue struct {
	Bundles []registry.Bundle
}

// benchmark bundles workflow using registry.Bundle{}
func BenchmarkBundlesByValue(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		bundles := make([]registry.Bundle, nBundles) // pre-allocate
		for j := 0; j < nBundles; j++ {
			bundle, err := testutils.NewBundle(benchmarkBundleName, bundlePath)
			require.NoError(b, err)
			bundles[j] = *bundle
		}
		require.Equal(b, len(bundles), nBundles)
		mb := &metaBundleByValue{Bundles: bundles}
		handleBundlesByValue(b, mb)
	}
}

func handleBundlesByValue(b *testing.B, mb *metaBundleByValue) {
	for _, bundle := range mb.Bundles {
		bundleName := getBundleNameByValue(bundle)
		require.Equal(b, bundleName, benchmarkBundleName)
	}
}

func getBundleNameByValue(bundle registry.Bundle) string {
	return bundle.Name
}
