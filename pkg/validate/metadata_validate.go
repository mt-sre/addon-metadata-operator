package validate

import (
	"github.com/mt-sre/addon-metadata-operator/pkg/validators/meta"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators/cross"
	
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"

)

func GetAllMetaValidators() []utils.Validator {
	return []utils.Validator{
		{
			Description: "Ensure defaultChannel is present in list of channels",
			Runner:      meta.ValidateDefaultChannel,
		},
		{
			Description: "Ensure `label` to follow the format api.openshift.com/addon-<operator-name>",
			Runner:      meta.ValidateAddonLabel,
		},
		{
			Description: "Some description about some cross validator",
			Runner:      cross.ValidateSomethingCrossBetweenImageSetAndMetadata, // cross validators can be separated into a different function as well like GetAllCrossValidators() []utils.Validator
		},
	}
}
