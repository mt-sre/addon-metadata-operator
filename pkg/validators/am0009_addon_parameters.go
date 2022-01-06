package validators

import (
	"fmt"
	"regexp"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

func init() {
	Registry.Add(AM0009)
}

var AM0009 = utils.Validator{
	Code:        "AM0009",
	Name:        "addon_parameters",
	Description: "Ensure `addOnParameters` section in the addon metadata is rightfully defined",
	Runner:      ValidateAddonParameters,
}

func ValidateAddonParameters(metabundle utils.MetaBundle) (bool, string, error) {
	addonParams := metabundle.AddonMeta.AddOnParameters
	if addonParams == nil {
		return true, "", nil
	}
	for _, param := range *addonParams {
		validation := param.Validation
		options := param.Options
		defaultValue := param.DefaultValue

		if validation != nil && options != nil {
			return false, "validation and options can't both be set", nil
		}

		if defaultValue != nil {
			if validation != nil {
				r, err := regexp.Compile(*validation)
				if err != nil {
					return false, "", fmt.Errorf("failed parse `validation` as regex: %w", err)
				}
				if !r.MatchString(*defaultValue) {
					return false, fmt.Sprintf("defaultValue %s didn't match its validation", *defaultValue), nil
				}
				return true, "", nil
			}

			if options != nil {
				for _, opt := range *options {
					if *defaultValue == opt.Value {
						return true, "", nil
					}
				}
				return false, fmt.Sprintf("defaultValue '%s' not found in `options`", *defaultValue), nil
			}
		}
	}
	return true, "", nil
}
