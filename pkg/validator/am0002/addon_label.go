package am0002

import (
	"context"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewAddonLabel)
}

const (
	code = 2
	name = "label_format"
	desc = "Validates whether label follows the format 'api.openshift.com/addon-<id>'"
)

func NewAddonLabel(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &AddonLabel{
		Base: base,
	}, nil
}

type AddonLabel struct {
	*validator.Base
}

func (a *AddonLabel) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	operatorId, label := mb.AddonMeta.ID, mb.AddonMeta.Label
	if label != "api.openshift.com/addon-"+operatorId {
		msg := fmt.Sprintf("addon label '%s' wasn't recognized to follow the 'api.openshift.com/addon-<id>' format", label)
		return a.Fail(msg)
	}

	return a.Success()
}
