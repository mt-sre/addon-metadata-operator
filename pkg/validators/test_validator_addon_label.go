package validators

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

type ValidatorAddonLabelTestBundle struct{}

func (val ValidatorAddonLabelTestBundle) Name() string {
	return "Addon Label Validator"
}

func (val ValidatorAddonLabelTestBundle) Run(mb utils.MetaBundle) (bool, error) {
	return ValidateAddonLabel(mb)
}

func (val ValidatorAddonLabelTestBundle) SucceedingCandidates() []utils.MetaBundle {
	return []utils.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:    "random-operator",
				Label: "api.openshift.com/addon-random-operator",
			},
		},
	}
}

func (val ValidatorAddonLabelTestBundle) FailingCandidates() []utils.MetaBundle {
	return []utils.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:    "random-operator",
				Label: "foo-bar",
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:    "random-operator",
				Label: "api.openshift.com/addon-random-operator-x",
			},
		},
	}
}
