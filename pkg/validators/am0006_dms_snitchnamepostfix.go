package validators

import (
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

func init() {
	Registry.Add(AM0006)
}

var AM0006 = types.NewValidator(
	"AM0006",
	types.ValidateFunc(ValidateDmsSnitchNamePostFix),
	types.ValidatorName("dms_snitchnamepostfix"),
	types.ValidatorDescription("Ensure `deadmanssnitch.snitchNamePostFix` doesn't begin with 'hive-'"),
)

func ValidateDmsSnitchNamePostFix(cfg types.ValidatorConfig, mb types.MetaBundle) types.ValidatorResult {
	dmsConf := mb.AddonMeta.DeadmansSnitch
	if dmsConf == nil || dmsConf.SnitchNamePostFix == nil {
		return Success()
	}
	if strings.HasPrefix(*dmsConf.SnitchNamePostFix, "hive-") {
		return Fail(fmt.Sprintf("`deadmanssnitch.snitchNamePostFix` in addon %s found to begin with 'hive-'", mb.AddonMeta.ID))
	}
	return Success()
}
