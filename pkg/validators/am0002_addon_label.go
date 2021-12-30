package validators

import (
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

func init() {
	Registry.Add(AM0002)
}

var AM0002 = types.Validator{
	Code:        "AM0002",
	Name:        "label_format",
	Description: "Validates whether label follows the format 'api.openshift.com/addon-<id>'",
	Runner:      validateAddonLabel,
}

func validateAddonLabel(mb types.MetaBundle) (bool, string, error) {
	operatorId, label := mb.AddonMeta.ID, mb.AddonMeta.Label
	if label != "api.openshift.com/addon-"+operatorId {
		return false, fmt.Sprintf("addon label '%s' wasn't recognized to follow the 'api.openshift.com/addon-<id>' format", label), nil
	}

	return true, "", nil
}
