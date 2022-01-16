package validators_test

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	mtsrev1 "github.com/mt-sre/addon-metadata-operator/pkg/mtsre/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
)

func init() {
	TestRegistry.Add(TestAM0006{})
}

type TestAM0006 struct{}

func (t TestAM0006) Name() string {
	return validators.AM0006.Name
}

func (t TestAM0006) Run(mb types.MetaBundle) types.ValidatorResult {
	return validators.AM0006.Validate(mb)
}

func (t TestAM0006) SucceedingCandidates() ([]types.MetaBundle, error) {
	res, err := testutils.DefaultSucceedingCandidates()
	if err != nil {
		return nil, err
	}

	moreSucceedingCandidates := []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "random-operator-0",
				DeadmansSnitch: &mtsrev1.DeadmansSnitch{
					SnitchNamePostFix: testutils.GetStringLiteralRef("addon-123"),
				},
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator-1",
				DeadmansSnitch: &mtsrev1.DeadmansSnitch{},
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "random-operator-2",
			},
		},
	}

	return append(res, moreSucceedingCandidates...), nil
}

func (t TestAM0006) FailingCandidates() ([]types.MetaBundle, error) {
	res := []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "random-operator",
				DeadmansSnitch: &mtsrev1.DeadmansSnitch{
					SnitchNamePostFix: testutils.GetStringLiteralRef("hive-123"),
				},
			},
		},
	}
	return res, nil
}
