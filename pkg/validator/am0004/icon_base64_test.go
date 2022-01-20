package am0004

import (
	_ "embed"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/stretchr/testify/require"
)

//go:embed valid_icon.b64
var validIcon string

func TestIconBase64Valid(t *testing.T) {
	t.Parallel()

	bundles, err := testutils.DefaultValidBundleMap()
	require.NoError(t, err)

	for name, bundle := range map[string]types.MetaBundle{
		"valid base64 encoded png": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:   "random-operator-1",
				Icon: validIcon,
			},
		},
	} {
		bundles[name] = bundle
	}

	tester := testutils.NewValidatorTester(t, NewIconBase64)
	tester.TestValidBundles(bundles)
}

func TestIconBase64Invalid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewIconBase64)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"no icon": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "random-operator-3",
			},
		},
		"invalid base64": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:   "random-operator-4",
				Icon: "not-base64",
			},
		},
		"invalid png": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:   "random-operator-4",
				Icon: "dGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIHRoZSBsYXp5IGRvZw==",
			},
		},
	})
}
