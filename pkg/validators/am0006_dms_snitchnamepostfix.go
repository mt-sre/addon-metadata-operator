package validators

import (
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

func init() {
	Registry.Add(AM0006)
}

var AM0006 = utils.Validator{
	Code:        "AM0006",
	Name:        "dms_snitchnamepostfix",
	Description: "Ensure `deadmanssnitch.snitchNamePostFix` doesn't begin with 'hive-'",
	Runner:      ValidateDmsSnitchNamePostFix,
}

func ValidateDmsSnitchNamePostFix(metabundle utils.MetaBundle) (bool, string, error) {
	dmsConf := metabundle.AddonMeta.DeadmansSnitch
	if dmsConf == nil || dmsConf.SnitchNamePostFix == nil {
		return true, "", nil
	}
	if strings.HasPrefix(*dmsConf.SnitchNamePostFix, "hive-") {
		return false, fmt.Sprintf("`deadmanssnitch.snitchNamePostFix` in addon %s found to begin with 'hive-'", metabundle.AddonMeta.ID), nil
	}
	return true, "", nil
}
