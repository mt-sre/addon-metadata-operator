package am0012

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
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
	invalidBundles, err := utils.ExtractAndParseAddons(
		"quay.io/osd-addons/rhods-index@sha256:487e106059aea611af377985e6f30d7879bc36c4a16fe0f70531b7c1befd4675",
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
