package validators

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

type ValidatorDefaultChannelTestBundle struct{}

func (val ValidatorDefaultChannelTestBundle) Name() string {
	return "Addon Default Channel Validator"
}

func (val ValidatorDefaultChannelTestBundle) Run(mb utils.MetaBundle) (bool, error) {
	return ValidateDefaultChannel(mb)
}

func (val ValidatorDefaultChannelTestBundle) SucceedingCandidates() []utils.MetaBundle {
	return []utils.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator",
				DefaultChannel: "alpha",
				Channels: []v1alpha1.Channel{
					{
						Name: "alpha",
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
				Channels: []v1alpha1.Channel{
					{
						Name: "alpha",
					},
					{
						Name: "beta",
					},
				},
			},
		},
	}
}

func (val ValidatorDefaultChannelTestBundle) FailingCandidates() []utils.MetaBundle {
	return []utils.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator",
				DefaultChannel: "alpha",
				Channels: []v1alpha1.Channel{
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
				Channels: []v1alpha1.Channel{
					{
						Name: "alpha",
					},
				},
			},
		},
	}
}
