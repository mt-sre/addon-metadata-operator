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
	Description: "Validates whether label follows the format 'api.openshift.com/addon-<id>'",
	Runner:      validateAddonLabel,
}

func validateAddonLabel(metabundle utils.MetaBundle) (bool, string, error) {
	operatorId, label := metabundle.AddonMeta.ID, metabundle.AddonMeta.Label
	if label != "api.openshift.com/addon-"+operatorId {
		return false, fmt.Sprintf("addon label '%s' wasn't recognized to follow the 'api.openshift.com/addon-<id>' format", label), nil
	}

	return true, "", nil
}
