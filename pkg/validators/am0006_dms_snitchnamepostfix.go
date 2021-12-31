package validators

import (
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

func init() {
	Registry.Add(AM0006)
}

var AM0006 = types.Validator{
	Code:        "AM0006",
	Name:        "dms_snitchnamepostfix",
	Description: "Ensure `deadmanssnitch.snitchNamePostFix` doesn't begin with 'hive-'",
	Runner:      ValidateDmsSnitchNamePostFix,
}

func ValidateDmsSnitchNamePostFix(mb types.MetaBundle) types.ValidatorResult {
	dmsConf := mb.AddonMeta.DeadmansSnitch
	if dmsConf == nil || dmsConf.SnitchNamePostFix == nil {
		return Success()
	}
	if strings.HasPrefix(*dmsConf.SnitchNamePostFix, "hive-") {
		return Fail(fmt.Sprintf("`deadmanssnitch.snitchNamePostFix` in addon %s found to begin with 'hive-'", mb.AddonMeta.ID))
	}
	return Success()
}
