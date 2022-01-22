package extractor

import (
	"context"
	"testing"

	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	bundleImage         string
	expectedPackageName string
	expectedCSVName     string
	expectedCSVVersion  string
}

func TestDefaultBundleExtractorImplements(t *testing.T) {
	require.Implements(t, new(BundleExtractor), &DefaultBundleExtractor{})
}

func TestExtractorInMemoryCacheJSONSnappyEncoder(t *testing.T) {
	cases := []testCase{
		// reference-addon:0.1.6
		{
			bundleImage:         "quay.io/osd-addons/reference-addon-bundle@sha256:a62fd3f3b55aa58c587f0b7630f5e70b123d036a1a04a1bd5a866b5c576a04f4",
			expectedPackageName: "reference-addon",
			expectedCSVName:     "reference-addon.v0.1.6",
			expectedCSVVersion:  "0.1.6",
		},
		// reference-addon:0.1.5
		{
			bundleImage:         "quay.io/osd-addons/reference-addon-bundle@sha256:29879d193bd8da42e7b6500252b4d21bef733666bd893de2a3f9b250e591658e",
			expectedPackageName: "reference-addon",
			expectedCSVName:     "reference-addon.v0.1.5",
			expectedCSVVersion:  "0.1.5",
		},
	}
	encoder := NewJSONSnappyEncoder()
	cache := NewBundleMemoryCache(encoder)

	// adding extra logging for easier test debugging
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	extractor := NewBundleExtractor(WithBundleCache(cache), WithBundleLog(log))

	for _, tc := range cases {
		tc := tc // pin
		t.Run(tc.bundleImage, func(t *testing.T) {
			t.Parallel()
			bundle, err := extractor.Extract(context.Background(), tc.bundleImage)
			require.NoError(t, err)
			require.NotNil(t, bundle)
			testBundleFields(t, bundle, tc)

			cachedBundle, err := cache.Get(tc.bundleImage)
			require.NoError(t, err)
			testBundleFields(t, cachedBundle, tc)

			// make sure we encode to same length as an equality check
			n1, err := encoder.Encode(bundle)
			require.NoError(t, err)
			n2, err := encoder.Encode(cachedBundle)
			require.NoError(t, err)
			require.Equal(t, len(n1), len(n2))
		})
	}
}

func testBundleFields(t *testing.T, b *registry.Bundle, tc testCase) {
	require.Equal(t, b.Name, tc.expectedPackageName)
	require.Equal(t, b.BundleImage, tc.bundleImage)
	require.NotNil(t, b.Annotations)
	require.Equal(t, b.Annotations.PackageName, tc.expectedPackageName)

	csv, err := b.ClusterServiceVersion()
	require.NoError(t, err)
	require.NotNil(t, csv)
	require.Equal(t, csv.Name, tc.expectedCSVName)

	version, err := csv.GetVersion()
	require.NoError(t, err)
	require.Equal(t, version, tc.expectedCSVVersion)
}
