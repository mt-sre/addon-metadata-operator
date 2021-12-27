package validators

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

func init() {
	TestRegistry.Add(TestAM0001{})
}

type TestAM0001 struct{}

func (t TestAM0001) Name() string {
	return AM0001.Name
}

func (t TestAM0001) Run(mb utils.MetaBundle) (bool, string, error) {
	return AM0001.Runner(mb)
}

func (t TestAM0001) SucceedingCandidates() []utils.MetaBundle {
	return testutils.DefaultSucceedingCandidates()
}

func (t TestAM0001) FailingCandidates() []utils.MetaBundle {
	return []utils.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator",
				DefaultChannel: "alpha",
				Channels: &[]v1alpha1.Channel{
					{
						Name: "beta",
					},
					{
						Name: "sigma",
					},
				},
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator",
				DefaultChannel: "beta",
				Channels: &[]v1alpha1.Channel{
					{
						Name: "alpha",
					},
				},
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator",
				DefaultChannel: "invalid",
			},
		},
	}
}
