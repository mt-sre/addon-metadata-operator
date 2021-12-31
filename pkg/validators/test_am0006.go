package validators

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	mtsrev1 "github.com/mt-sre/addon-metadata-operator/pkg/mtsre/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

func init() {
	TestRegistry.Add(TestAM0006{})
}

type TestAM0006 struct{}

func (t TestAM0006) Name() string {
	return AM0006.Name
}

func (t TestAM0006) Run(mb types.MetaBundle) types.ValidatorResult {
	return AM0006.Runner(mb)
}

func (t TestAM0006) SucceedingCandidates() []types.MetaBundle {
	res := testutils.DefaultSucceedingCandidates()

	moreSucceedingCandidates := []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "random-operator-0",
				DeadmansSnitch: &mtsrev1.DeadmansSnitch{
					SnitchNamePostFix: StringToPtr("addon-123"),
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

	return append(res, moreSucceedingCandidates...)
}

func (t TestAM0006) FailingCandidates() []types.MetaBundle {
	return []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "random-operator",
				DeadmansSnitch: &mtsrev1.DeadmansSnitch{
					SnitchNamePostFix: StringToPtr("hive-123"),
				},
			},
		},
	}
}
