package am0013

import (
	"context"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewAddonRequirements)
}

const (
	code = 13
	name = "addon_requirements"
	desc = "Ensure `addOnRequirements` section in the addon metadata is rightfully defined"
)

func NewAddonRequirements(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &AddonRequirements{
		Base: base,
	}, nil
}

type AddonRequirements struct {
	*validator.Base
}

func (a *AddonRequirements) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	requirements := mb.AddonMeta.AddOnRequirements

	if requirements == nil {
		return a.Success()
	}

	var msgs []string

	for _, r := range *requirements {
		if len(r.Data) == 0 {
			msgs = append(msgs, fmt.Sprintf("requirement %q has no data", r.ID))
		}
	}

	if len(msgs) == 0 {
		return a.Success()
	}

	return a.Fail(msgs...)
}
