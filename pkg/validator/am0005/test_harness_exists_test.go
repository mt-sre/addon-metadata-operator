package am0005

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/stretchr/testify/require"
)

func TestTestHarnessExistsValid(t *testing.T) {
	t.Parallel()

	bundles, err := testutils.DefaultValidBundleMap()
	require.NoError(t, err)

	for name, bundle := range map[string]types.MetaBundle{
		"untagged harness image": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "quay.io/miwilson/addon-samples",
			},
		},
		"tagged harness image": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "quay.io/asnaraya/reference-addon-test-harness:fix",
			},
		},
		"hashed harness image": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "quay.io/asnaraya/reference-addon-test-harness@sha256:bdc32a600202d36fec4524dbec177e9313ef82ad4bda5bd24d4b75236ca8a482",
			},
		},
	} {
		bundles[name] = bundle
	}

	tester := testutils.NewValidatorTester(t, NewTestHarnessExists)
	tester.TestValidBundles(bundles)
}

func TestTestHarnessExistsInvalid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewTestHarnessExists)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"invalid url": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "abcd",
			},
		},
		"non existent harness": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "quay.io/asnaraya/reference-addon-test-harness:404",
			},
		},
		"non quay hosted image": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "https://docker.io/ashishmax31/addon-operator-bundle:0.1.0-cb328d9",
			},
		},
	})
}
