package am0006

import (
	"context"
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewDMSSnitchNamePostFix)
}

const (
	code = 6
	name = "dms_snitchnamepostfix"
	desc = "Ensure `deadmanssnitch.snitchNamePostFix` doesn't begin with 'hive-'"
)

func NewDMSSnitchNamePostFix(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)

	return &DMSSnitchNamePostFix{
		Base: base,
	}, err
}

type DMSSnitchNamePostFix struct {
	*validator.Base
}

func (d *DMSSnitchNamePostFix) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	dmsConf := mb.AddonMeta.DeadmansSnitch
	if dmsConf == nil || dmsConf.SnitchNamePostFix == nil {
		return d.Success()
	}
	if strings.HasPrefix(*dmsConf.SnitchNamePostFix, "hive-") {
		return d.Fail(fmt.Sprintf("`deadmanssnitch.snitchNamePostFix` in addon %s found to begin with 'hive-'", mb.AddonMeta.ID))
	}
	return d.Success()
}
