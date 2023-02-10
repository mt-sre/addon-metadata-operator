package am0012

import (
	"context"
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
	ctx := context.Background()

	extractor := extractor.New()
	invalidBundles, err := extractor.ExtractBundles(
		ctx,
		"quay.io/osd-addons/rhods-index@sha256:cd579e328ecce9141793c8df4dcb6675fac93a61035003a2f0820e37db616f74",
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
