package extractor

import (
	"context"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/operator"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultBundleExtractorImplements(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(BundleExtractor), &DefaultBundleExtractor{})
}

func TestDefaultBundleExtractor(t *testing.T) {
	t.Parallel()

	cache := NewBundleCacheImpl()

	// adding extra logging for easier test debugging
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	extractor := NewBundleExtractor(WithBundleCache(cache), WithBundleLog(log))

	for name, tc := range map[string]testCase{
		"reference-addon:0.1.6": {
			BundleImage:         "quay.io/osd-addons/reference-addon-bundle@sha256:a62fd3f3b55aa58c587f0b7630f5e70b123d036a1a04a1bd5a866b5c576a04f4",
			ExpectedPackageName: "reference-addon",
			ExpectedCSVName:     "reference-addon.v0.1.6",
			ExpectedCSVVersion:  "0.1.6",
		},
		"reference-addon:0.1.5": {
			BundleImage:         "quay.io/osd-addons/reference-addon-bundle@sha256:29879d193bd8da42e7b6500252b4d21bef733666bd893de2a3f9b250e591658e",
			ExpectedPackageName: "reference-addon",
			ExpectedCSVName:     "reference-addon.v0.1.5",
			ExpectedCSVVersion:  "0.1.5",
		},
	} {
		tc := tc // pin

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bundle, err := extractor.Extract(context.Background(), tc.BundleImage)
			require.NoError(t, err)

			require.NotNil(t, bundle)
			tc.AssertExpectations(t, bundle)

			cachedBundle, err := cache.GetBundle(tc.BundleImage)
			require.NoError(t, err)

			tc.AssertExpectations(t, *cachedBundle)
		})
	}
}

type testCase struct {
	BundleImage         string
	ExpectedPackageName string
	ExpectedCSVName     string
	ExpectedCSVVersion  string
}

func (tc testCase) AssertExpectations(t *testing.T, b operator.Bundle) {
	t.Helper()

	assert.Equal(t, b.Name, tc.ExpectedPackageName)
	assert.Equal(t, b.BundleImage, tc.BundleImage)
	assert.NotNil(t, b.Annotations)
	assert.Equal(t, b.Annotations.PackageName, tc.ExpectedPackageName)

	assert.Equal(t, b.ClusterServiceVersion.Name, tc.ExpectedCSVName)

	assert.Equal(t, b.Version, tc.ExpectedCSVVersion)
}
