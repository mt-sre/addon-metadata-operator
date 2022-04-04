package am0002

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
)

func TestAddonLabelValid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewAddonLabel)

	tester.TestValidBundles(map[string]types.MetaBundle{
		"properly prefixed label": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:    "random-operator",
				Label: "api.openshift.com/addon-random-operator",
			},
		},
	})

}

func TestAddonLabelInvalid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewAddonLabel)

	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"no prefix": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:    "random-operator",
				Label: "foo-bar",
			},
		},
		"non matching addon id": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:    "random-operator",
				Label: "api.openshift.com/addon-random-operator-x",
			},
		},
	})
}
