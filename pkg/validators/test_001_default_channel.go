<<<<<<< HEAD
package validators

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

// check interface implemented
var _ = utils.ValidatorTest(ValidatorTest001DefaultChannel{})

type ValidatorTest001DefaultChannel struct{}

func (v ValidatorTest001DefaultChannel) Name() string {
	return "Addon Default Channel Validator"
}

func (v ValidatorTest001DefaultChannel) Run(mb utils.MetaBundle) (bool, string, error) {
	return Validate001DefaultChannel(mb)
}

func (v ValidatorTest001DefaultChannel) SucceedingCandidates() []utils.MetaBundle {
	return testutils.DefaultSucceedingCandidates()
}

// not implemented
func (v ValidatorTest001DefaultChannel) FailingCandidates() []utils.MetaBundle {
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
||||||| parent of e58bb06 (Refactors.)
=======
package validators

import (
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

// check interface implemented
var _ = utils.ValidatorTest(ValidatorTest001DefaultChannel{})

type ValidatorTest001DefaultChannel struct{}

func (v ValidatorTest001DefaultChannel) Name() string {
	return "Addon Default Channel Validator"
}

func (v ValidatorTest001DefaultChannel) Run(mb utils.MetaBundle) (bool, string, error) {
	return Validate001DefaultChannel(mb)
}

func (v ValidatorTest001DefaultChannel) SucceedingCandidates() []utils.MetaBundle {
	return testutils.DefaultSucceedingCandidates()
}

// not implemented
func (v ValidatorTest001DefaultChannel) FailingCandidates() []utils.MetaBundle {
	return []utils.MetaBundle{}
}
>>>>>>> e58bb06 (Refactors.)
