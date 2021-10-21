package validators

import (
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

// TODO: preferable, have return type as a (bool, error) instead
func ValidateAddonLabel(metabundle *utils.MetaBundle) (bool, error) {
	operatorName, label := metabundle.AddonMeta.OperatorName, metabundle.AddonMeta.Label
	if label != "api.openshift.com/addon-"+operatorName {
		return false, nil
	}

	return true, nil
}

func someOtherUtilFunction() {

}
