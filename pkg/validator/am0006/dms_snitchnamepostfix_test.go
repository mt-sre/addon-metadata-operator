package am0006

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	utils "github.com/mt-sre/addon-metadata-operator/internal/testutils"
	mtsrev1 "github.com/mt-sre/addon-metadata-operator/pkg/mtsre/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/stretchr/testify/require"
)

func TestDMSSnitchNamePostFixValid(t *testing.T) {
	t.Parallel()

	bundles, err := testutils.DefaultValidBundleMap()
	require.NoError(t, err)

	for name, bundle := range map[string]types.MetaBundle{
		"postFix defined without hive prefix": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "random-operator-0",
				DeadmansSnitch: &mtsrev1.DeadmansSnitch{
					SnitchNamePostFix: utils.GetStringLiteralRef("addon-123"),
				},
			},
		},
		"DMS defined with no postFix": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator-1",
				DeadmansSnitch: &mtsrev1.DeadmansSnitch{},
			},
		},
		"DMS not defined": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "random-operator-2",
			},
		},
	} {
		bundles[name] = bundle
	}

	tester := testutils.NewValidatorTester(t, NewDMSSnitchNamePostFix)
	tester.TestValidBundles(bundles)
}

func TestDMSSnitchNamePostFixInvalid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewDMSSnitchNamePostFix)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"hive prefixed postFix": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "random-operator",
				DeadmansSnitch: &mtsrev1.DeadmansSnitch{
					SnitchNamePostFix: utils.GetStringLiteralRef("hive-123"),
				},
			},
		},
	})
}
