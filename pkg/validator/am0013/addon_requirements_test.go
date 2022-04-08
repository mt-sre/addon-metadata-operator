package am0013

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	ocmv1 "github.com/mt-sre/addon-metadata-operator/pkg/ocm/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	utils "github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestAddonParametersValid(t *testing.T) {
	t.Parallel()

	bundles, err := utils.DefaultValidBundleMap()
	require.NoError(t, err)

	for name, bundle := range map[string]types.MetaBundle{
		"no addon requirements": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{},
		},
		"single addon requirement with data": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AddOnRequirements: &[]ocmv1.AddOnRequirement{
					{
						ID: "has data",
						Data: ocmv1.AddOnRequirementData{
							"key": apiextensionsv1.JSON{},
						},
					},
				},
			},
		},
		"multiple addon requirement with data": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AddOnRequirements: &[]ocmv1.AddOnRequirement{
					{
						ID: "has data",
						Data: ocmv1.AddOnRequirementData{
							"key": apiextensionsv1.JSON{},
						},
					},
					{
						ID: "has more data",
						Data: ocmv1.AddOnRequirementData{
							"key": apiextensionsv1.JSON{},
						},
					},
				},
			},
		},
	} {
		bundles[name] = bundle
	}

	tester := utils.NewValidatorTester(t, NewAddonRequirements)
	tester.TestValidBundles(bundles)
}

func TestAddonParametersInvalid(t *testing.T) {
	t.Parallel()

	tester := utils.NewValidatorTester(t, NewAddonRequirements)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"single addon requirement without data": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AddOnRequirements: &[]ocmv1.AddOnRequirement{
					{
						ID:   "has no data",
						Data: ocmv1.AddOnRequirementData{},
					},
				},
			},
		},
		"multiple addon requirement without data": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AddOnRequirements: &[]ocmv1.AddOnRequirement{
					{
						ID:   "has no data",
						Data: ocmv1.AddOnRequirementData{},
					},
					{
						ID:   "also has no data",
						Data: ocmv1.AddOnRequirementData{},
					},
				},
			},
		},
		"multiple addon requirement with and without data": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AddOnRequirements: &[]ocmv1.AddOnRequirement{
					{
						ID:   "has no data",
						Data: ocmv1.AddOnRequirementData{},
					},
					{
						ID: "has data",
						Data: ocmv1.AddOnRequirementData{
							"key": apiextensionsv1.JSON{},
						},
					},
				},
			},
		},
	})
}
