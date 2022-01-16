package validators

import (
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

func init() {
	Registry.Add(AM0002)
}

var AM0002 = types.NewValidator(
	"AM0002",
	types.ValidateFunc(validateAddonLabel),
	types.ValidatorName("label_format"),
	types.ValidatorDescription("Validates whether label follows the format 'api.openshift.com/addon-<id>'"),
)

func validateAddonLabel(cfg types.ValidatorConfig, mb types.MetaBundle) types.ValidatorResult {
	operatorId, label := mb.AddonMeta.ID, mb.AddonMeta.Label
	if label != "api.openshift.com/addon-"+operatorId {
		msg := fmt.Sprintf("addon label '%s' wasn't recognized to follow the 'api.openshift.com/addon-<id>' format", label)
		return Fail(msg)
	}

	return Success()
}
