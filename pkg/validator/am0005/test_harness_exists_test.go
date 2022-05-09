package am0005

import (
	"context"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	imageparser "github.com/novln/docker-parser"
	"github.com/stretchr/testify/require"
)

func TestTestHarnessExistsValid(t *testing.T) {
	t.Parallel()

	quay := testutils.NewMockQuayClient()
	quay.
		On("HasReference",
			context.Background(),
			getRef(t, "quay.io/miwilson/addon-samples"),
		).
		Return(true, nil).
		On("HasReference",
			context.Background(),
			getRef(t, "quay.io/valid/no-tag"),
		).
		Return(true, nil).
		On("HasReference",
			context.Background(),
			getRef(t, "quay.io/valid/tag:tag"),
		).
		Return(true, nil).
		On("HasReference",
			context.Background(),
			getRef(t, "quay.io/valid/hash@sha256:bdc32a600202d36fec4524dbec177e9313ef82ad4bda5bd24d4b75236ca8a482"),
		).
		Return(true, nil)

	bundles, err := testutils.DefaultValidBundleMap()
	require.NoError(t, err)

	for name, bundle := range map[string]types.MetaBundle{
		"untagged harness image": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "quay.io/valid/no-tag",
			},
		},
		"tagged harness image": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "quay.io/valid/tag:tag",
			},
		},
		"hashed harness image": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "quay.io/valid/hash@sha256:bdc32a600202d36fec4524dbec177e9313ef82ad4bda5bd24d4b75236ca8a482",
			},
		},
	} {
		bundles[name] = bundle
	}

	tester := testutils.NewValidatorTester(t,
		NewTestHarnessExists,
		testutils.ValidatorTesterQuayClient(quay),
	)
	tester.TestValidBundles(bundles)
}

func TestTestHarnessExistsInvalid(t *testing.T) {
	t.Parallel()

	quay := testutils.NewMockQuayClient()
	quay.
		On("HasReference",
			context.Background(),
			getRef(t, "abcd"),
		).
		Return(false, nil).
		On("HasReference",
			context.Background(),
			getRef(t, "quay.io/asnaraya/reference-addon-test-harness:404"),
		).
		Return(false, nil).
		On("HasReference",
			context.Background(),
			getRef(t, "https://docker.io/ashishmax31/addon-operator-bundle:0.1.0-cb328d9"),
		).
		Return(false, nil)

	tester := testutils.NewValidatorTester(t,
		NewTestHarnessExists,
		testutils.ValidatorTesterQuayClient(quay),
	)
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

func getRef(t *testing.T, image string) *imageparser.Reference {
	t.Helper()

	ref, err := imageparser.Parse(image)
	require.NoError(t, err)

	return ref
}
