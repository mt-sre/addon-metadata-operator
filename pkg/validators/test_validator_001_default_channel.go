package validators

import (
	"log"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

type Validator001DefaultChannel struct{}

func (v Validator001DefaultChannel) Name() string {
	return "Addon Default Channel Validator"
}

func (v Validator001DefaultChannel) Run(mb utils.MetaBundle) (bool, string, error) {
	return ValidateDefaultChannel(mb)
}

func (v Validator001DefaultChannel) SucceedingCandidates() []utils.MetaBundle {
	res, err := testutils.DefaultSucceedingCandidates()
	if err != nil {
		log.Fatalf("Could not load default succeeding candidates, got %v. Exiting.", err)
	}
	return res
}

// TODO - get succeeding candidates and modify them???
func (v Validator001DefaultChannel) FailingCandidates() []utils.MetaBundle {
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
	}
}
