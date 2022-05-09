package am0011

import (
	"context"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/stretchr/testify/require"
)

func TestSKURuleExistsValid(t *testing.T) {
	t.Parallel()

	ocm := testutils.NewMockOCMClient()
	ocm.
		On("QuotaRuleExists", context.Background(), "addon-reference-addon").
		Return(true, nil).
		On("QuotaRuleExists", context.Background(), "addon-successful-candidate").
		Return(true, nil).
		On("QuotaRuleExists", context.Background(), "addon-zero-quota-candidate").
		Return(true, nil)

	bundles, err := testutils.DefaultValidBundleMap()
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

	tester := testutils.NewValidatorTester(
		t, NewOCMSKURuleExists,
		testutils.ValidatorTesterOCMClient(ocm),
	)
	tester.TestValidBundles(bundles)
}

func TestSKURuleExistsInvalid(t *testing.T) {
	t.Parallel()

	ocm := testutils.NewMockOCMClient()
	ocm.
		On("QuotaRuleExists", context.Background(), "addon-failing-candidate").
		Return(false, nil)

	tester := testutils.NewValidatorTester(
		t, NewOCMSKURuleExists,
		testutils.ValidatorTesterOCMClient(ocm),
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
