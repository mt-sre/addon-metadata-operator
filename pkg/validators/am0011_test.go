package validators_test

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
)

func init() {
	TestRegistry.Add(TestAM0011{})
}

type TestAM0011 struct{}

func (val TestAM0011) Name() string {
	return validators.AM0011.Name
}

func (val TestAM0011) Run(mb types.MetaBundle) types.ValidatorResult {
	client := testutils.NewMockOCMClient(
		testutils.MockOCMClientValidQuotaNames(
			"addon-reference-addon",
			"addon-successful-candidate",
			"addon-zero-quota-candidate",
		),
	)

	validateFunc := validators.GenerateOCMSKUValidator(client)

	return validateFunc(mb)
}

func (val TestAM0011) SucceedingCandidates() ([]types.MetaBundle, error) {
	candidates, err := testutils.DefaultSucceedingCandidates()
	if err != nil {
		return nil, err
	}

	return append(candidates,
		types.MetaBundle{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OcmQuotaName: "addon-successful-candidate",
				OcmQuotaCost: 1,
			},
		}, types.MetaBundle{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OcmQuotaName: "addon-zero-quota-candidate",
				OcmQuotaCost: 0,
			},
		}), nil
}

func (val TestAM0011) FailingCandidates() ([]types.MetaBundle, error) {
	return []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OcmQuotaName: "addon-failing-candidate",
				OcmQuotaCost: 1,
			},
		}}, nil
}
