package validators

import (
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

// TODO: preferable, have return type as a (bool, error) instead
func ValidateAddonLabel(metabundle *utils.MetaBundle) (bool, error) {
	operatorId, label := metabundle.AddonMeta.ID, metabundle.AddonMeta.Label
	if label != "api.openshift.com/addon-"+operatorId {
		return false, nil
	}

	return true, nil
}

func someOtherUtilFunction() {

}
