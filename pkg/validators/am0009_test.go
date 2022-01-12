package validators_test

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	ocmv1 "github.com/mt-sre/addon-metadata-operator/pkg/ocm/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
)

func init() {
	TestRegistry.Add(TestAM0009{})
}

type TestAM0009 struct{}

func (val TestAM0009) Name() string {
	return validators.AM0009.Name
}

func (val TestAM0009) Run(mb types.MetaBundle) types.ValidatorResult {
	return validators.AM0009.Runner(mb)
}

func (val TestAM0009) SucceedingCandidates() ([]types.MetaBundle, error) {
	res, err := testutils.DefaultSucceedingCandidates()
	if err != nil {
		return nil, err
	}

	moreCandidates := []types.MetaBundle{
		{
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
		{
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
		{
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
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "addon-parameters-not-specified",
			},
		},
	}

	return append(res, moreCandidates...), nil
}

func (val TestAM0009) FailingCandidates() ([]types.MetaBundle, error) {
	res := []types.MetaBundle{
		{
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
		{
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
		{
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
	}
	return res, nil
}
