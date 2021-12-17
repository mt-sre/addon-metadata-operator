package validators

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	mtsrev1 "github.com/mt-sre/addon-metadata-operator/pkg/mtsre/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

func init() {
	TestRegistry.Add(TestAM0006{})
}

type TestAM0006 struct{}

func (t TestAM0006) Name() string {
	return AM0006.Name
}

func (t TestAM0006) Run(mb utils.MetaBundle) (bool, string, error) {
	return AM0006.Runner(mb)
}

func (t TestAM0006) SucceedingCandidates() []utils.MetaBundle {
	res := testutils.DefaultSucceedingCandidates()

	moreSucceedingCandidates := []utils.MetaBundle{
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

func (t TestAM0006) FailingCandidates() []utils.MetaBundle {
	return []utils.MetaBundle{
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
