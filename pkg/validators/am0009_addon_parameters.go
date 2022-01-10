package validators

import (
	"fmt"
	"regexp"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

func init() {
	Registry.Add(AM0009)
}

var AM0009 = types.Validator{
	Code:        "AM0009",
	Name:        "addon_parameters",
	Description: "Ensure `addOnParameters` section in the addon metadata is rightfully defined",
	Runner:      ValidateAddonParameters,
}

func ValidateAddonParameters(mb types.MetaBundle) types.ValidatorResult {
	addonParams := mb.AddonMeta.AddOnParameters
	if addonParams == nil {
		return Success()
	}
	for _, param := range *addonParams {
		validation := param.Validation
		options := param.Options
		defaultValue := param.DefaultValue

		if validation != nil && options != nil {
			return Fail("validation and options can't both be set")
		}

		if defaultValue != nil {
			if validation != nil {
				r, err := regexp.Compile(*validation)
				if err != nil {
					return Error(fmt.Errorf("failed parse `validation` as regex: %w", err))
				}
				if !r.MatchString(*defaultValue) {
					return Fail(fmt.Sprintf("defaultValue %s didn't match its validation", *defaultValue))
				}
				return Success()
			}

			if options != nil {
				for _, opt := range *options {
					if *defaultValue == opt.Value {
						return Success()
					}
				}
				return Fail(fmt.Sprintf("defaultValue '%s' not found in `options`", *defaultValue))
			}
		}
	}
	return Success()
}
