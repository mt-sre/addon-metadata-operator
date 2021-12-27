package validators

import (
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

func init() {
	Registry.Add(AM0002)
}

var AM0002 = utils.Validator{
	Code:        "AM0002",
	Name:        "label_format",
	Description: "Ensure defaultChannel is present in list of channels",
	Runner:      ValidateAddonLabel,
}

// ValidateAddonLabel validates whether the 'label' field under an addon.yaml follows the format 'api.openshift.com/addon-<id>'
// TODO - remove this validator once we do field level validation
func ValidateAddonLabel(metabundle utils.MetaBundle) (bool, string, error) {
	operatorId, label := metabundle.AddonMeta.ID, metabundle.AddonMeta.Label
	if label != "api.openshift.com/addon-"+operatorId {
		return false, fmt.Sprintf("addon label '%s' wasn't recognized to follow the 'api.openshift.com/addon-<id>' format", label), nil
	}

	return true, "", nil
}
