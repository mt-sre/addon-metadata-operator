package am0012

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/extractor"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/stretchr/testify/require"
)

func TestCSVRBAC(t *testing.T) {
	bundles, err := testutils.DefaultValidBundleMap()
	require.NoError(t, err)

	tester := testutils.NewValidatorTester(t, NewCSVRBAC)
	tester.TestValidBundles(bundles)
}

func TestCSVRBACInvalid(t *testing.T) {
	extractor := extractor.New()
	invalidBundles, err := extractor.ExtractBundles(
		"quay.io/osd-addons/rhods-index@sha256:94c934b33096b057c07a4ffb5ae59a6f0c7641fe45d71fc1283182f6c01a8ef3",
		"rhods-operator",
	)
	require.NoError(t, err)

	tester := testutils.NewValidatorTester(t, NewCSVRBAC)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"non existing quota name": {
			Bundles: invalidBundles,
		},
	})

}
