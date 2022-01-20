package am0009

import (
	"context"
	"fmt"
	"regexp"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewAddonParameters)
}

const (
	code = 9
	name = "addon_parameters"
	desc = "Ensure `addOnParameters` section in the addon metadata is rightfully defined"
)

func NewAddonParameters(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &AddonParamters{
		Base: base,
	}, nil
}

type AddonParamters struct {
	*validator.Base
}

func (a *AddonParamters) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	addonParams := mb.AddonMeta.AddOnParameters
	if addonParams == nil {
		return a.Success()
	}
	for _, param := range *addonParams {
		validation := param.Validation
		options := param.Options
		defaultValue := param.DefaultValue

		if validation != nil && options != nil {
			return a.Fail("validation and options can't both be set")
		}

		if defaultValue != nil {
			if validation != nil {
				r, err := regexp.Compile(*validation)
				if err != nil {
					return a.Error(fmt.Errorf("failed parse `validation` as regex: %w", err))
				}
				if !r.MatchString(*defaultValue) {
					msg := fmt.Sprintf("defaultValue %s didn't match its validation", *defaultValue)
					if param.ValidationErrMsg != nil {
						return a.Fail(fmt.Sprintf("%s: %s", msg, *param.ValidationErrMsg))
					}
					return a.Fail(msg)
				}
				return a.Success()
			}

			if options != nil {
				for _, opt := range *options {
					if *defaultValue == opt.Value {
						return a.Success()
					}
				}
				return a.Fail(fmt.Sprintf("defaultValue '%s' not found in `options`", *defaultValue))
			}
		}
	}
	return a.Success()
}
