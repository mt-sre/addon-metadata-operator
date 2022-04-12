package extractor

import (
	"context"
	"fmt"
	"testing"

	"github.com/mt-sre/regtest"
	dockerparser "github.com/novln/docker-parser"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	Image               string
	TarFile             string
	expectedPackageName string
	expectedCSVName     string
	expectedCSVVersion  string
}

func TestDefaultBundleExtractorImplements(t *testing.T) {
	require.Implements(t, new(BundleExtractor), &DefaultBundleExtractor{})
}

func TestExtractorInMemoryCacheJSONSnappyEncoder(t *testing.T) {
	cases := []testCase{
		{
			Image:               "localhost/reference-addon-bundle:v0.1.6",
			TarFile:             "./testdata/reference-addon-bundle-v0.1.6.tar.gz",
			expectedPackageName: "reference-addon",
			expectedCSVName:     "reference-addon.v0.1.6",
			expectedCSVVersion:  "0.1.6",
		},
		{
			Image:               "localhost/reference-addon-bundle:v0.1.5",
			TarFile:             "./testdata/reference-addon-bundle-v0.1.5.tar.gz",
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
	extractor := NewBundleExtractor(WithBundleCache(cache), WithBundleLog(log), WithBundleSkipTLS)

	for _, tc := range cases {
		tc := tc // pin
		t.Run(tc.Image, func(t *testing.T) {
			t.Parallel()

			reg := regtest.StartRegistry(t, regtest.WithImage("quay.io/openshifttest/registry:2"))

			defer func() { _ = reg.Stop() }()

			require.NoError(t, reg.Load(tc.Image, tc.TarFile))

			fullRef := fmt.Sprintf("%s/%s", reg.Host(), getImageName(t, tc.Image))

			bundle, err := extractor.Extract(context.Background(), fullRef)
			require.NoError(t, err)
			require.NotNil(t, bundle)
			testBundleFields(t, bundle, tc)

			cachedBundle, err := cache.Get(fullRef)
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
	t.Helper()

	require.Equal(t, tc.expectedPackageName, b.Name)
	require.Equal(t, getImageName(t, tc.Image), getImageName(t, b.BundleImage))
	require.NotNil(t, b.Annotations)
	require.Equal(t, tc.expectedPackageName, b.Annotations.PackageName)

	csv, err := b.ClusterServiceVersion()
	require.NoError(t, err)
	require.NotNil(t, csv)
	require.Equal(t, tc.expectedCSVName, csv.Name)

	version, err := csv.GetVersion()
	require.NoError(t, err)
	require.Equal(t, tc.expectedCSVVersion, version)
}

func getImageName(t *testing.T, ref string) string {
	t.Helper()

	parsed, err := dockerparser.Parse(ref)
	require.NoError(t, err)

	return parsed.Name()
}
