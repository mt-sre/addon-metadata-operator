package am0011

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	utils "github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/stretchr/testify/require"
)

func TestSKURuleExistsValid(t *testing.T) {
	t.Parallel()

	bundles, err := utils.DefaultValidBundleMap()
	require.NoError(t, err)

	for name, bundle := range map[string]types.MetaBundle{
		"existing quota name": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OcmQuotaName: "addon-successful-candidate",
				OcmQuotaCost: 1,
			},
		},
		"zero quota": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OcmQuotaName: "addon-zero-quota-candidate",
				OcmQuotaCost: 0,
			},
		},
	} {
		bundles[name] = bundle
	}

	tester := utils.NewValidatorTester(
		t, NewOCMSKURuleExists,
		utils.ValidatorTesterOCMClient(testutils.NewMockOCMClient(
			testutils.MockOCMClientValidQuotaNames(
				"addon-reference-addon",
				"addon-successful-candidate",
				"addon-zero-quota-candidate",
			),
		)),
	)
	tester.TestValidBundles(bundles)
}

func TestSKURuleExistsInvalid(t *testing.T) {
	t.Parallel()

	tester := utils.NewValidatorTester(
		t, NewOCMSKURuleExists,
		utils.ValidatorTesterOCMClient(testutils.NewMockOCMClient(
			testutils.MockOCMClientValidQuotaNames(
				"addon-reference-addon",
				"addon-successful-candidate",
				"addon-zero-quota-candidate",
			),
		)),
	)

	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"non existing quota name": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OcmQuotaName: "addon-failing-candidate",
				OcmQuotaCost: 1,
			},
		},
	})
}
