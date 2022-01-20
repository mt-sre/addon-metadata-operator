package am0009

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	ocmv1 "github.com/mt-sre/addon-metadata-operator/pkg/ocm/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	utils "github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/stretchr/testify/require"
)

func TestAddonParametersValid(t *testing.T) {
	t.Parallel()

	bundles, err := utils.DefaultValidBundleMap()
	require.NoError(t, err)

	for name, bundle := range map[string]types.MetaBundle{
		"options set without validation": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "default-value-in-options",
				AddOnParameters: &[]ocmv1.AddOnParameter{
					{
						ID:   "size",
						Name: "Managed StorageCluster size",
						Options: &[]ocmv1.AddOnParameterOption{
							{
								Name:  "1 TiB",
								Value: "1",
							},
							{
								Name:  "4 TiB",
								Value: "4",
							},
						},
						DefaultValue: testutils.GetStringLiteralRef("1"),
					},
				},
			},
		},
		"validation set with matching default": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "default-value-passing-validation",
				AddOnParameters: &[]ocmv1.AddOnParameter{
					{
						ID:           "notification-email-0",
						Name:         "An email address for storage lifecycle notifications.",
						Validation:   testutils.GetStringLiteralRef("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])"),
						DefaultValue: testutils.GetStringLiteralRef("something@something.com"),
					},
				},
			},
		},
		"validation set with no default": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "default-value-not-specified",
				AddOnParameters: &[]ocmv1.AddOnParameter{
					{
						ID:         "notification-email-1",
						Name:       "An email address for storage lifecycle notifications.",
						Validation: testutils.GetStringLiteralRef("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])"),
					},
				},
			},
		},
		"no addon parameters": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "addon-parameters-not-specified",
			},
		},
	} {
		bundles[name] = bundle
	}

	tester := utils.NewValidatorTester(t, NewAddonParameters)
	tester.TestValidBundles(bundles)
}

func TestAddonParametersInvalid(t *testing.T) {
	t.Parallel()

	tester := utils.NewValidatorTester(t, NewAddonParameters)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"both options and validation defined": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "both-options-and-validation",
				AddOnParameters: &[]ocmv1.AddOnParameter{
					{
						ID:         "notification-email-0",
						Name:       "An email address for storage lifecycle notifications.",
						Validation: testutils.GetStringLiteralRef("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])"),
						Options: &[]ocmv1.AddOnParameterOption{
							{
								Name:  "1 TiB",
								Value: "1",
							},
							{
								Name:  "4 TiB",
								Value: "4",
							},
						},
					},
				},
			},
		},
		"default value does not match validation": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "default-value-failing-validation",
				AddOnParameters: &[]ocmv1.AddOnParameter{
					{
						ID:           "notification-email-1",
						Name:         "An email address for storage lifecycle notifications.",
						Validation:   testutils.GetStringLiteralRef("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\\])"),
						DefaultValue: testutils.GetStringLiteralRef("a-non-email"),
					},
				},
			},
		},
		"default value defined outside of options array": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "default-value-not-in-options",
				AddOnParameters: &[]ocmv1.AddOnParameter{
					{
						ID:   "size",
						Name: "Managed cluster size",
						Options: &[]ocmv1.AddOnParameterOption{
							{
								Name:  "1 TiB",
								Value: "1",
							},
							{
								Name:  "4 TiB",
								Value: "4",
							},
						},
						DefaultValue: testutils.GetStringLiteralRef("not-in-options"),
					},
				},
			},
		},
	})
}
