package validate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/validators"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

func Validate(mb utils.MetaBundle, filter *validatorsFilter) (bool, []error) {
	errs := []error{}
	allSuccess := true

	printMetaHeading()

	for validatorName, validator := range filter.GetValidators() {
		fmt.Printf("\r%s\t\t", validator.Description)
		success, failureMsg, err := validator.Runner(mb)
		if err != nil {
			errs = append(errs, err)
			printErrorMessage(validator.Description)
		} else if !success {
			printFailureMessage(fmt.Sprintf("%v: %v.", validatorName, failureMsg))
			allSuccess = false
		} else {
			printSuccessMessage(validator.Description)
		}
		fmt.Println()
	}
	return allSuccess, errs
}

// name formatting rule: [0-9]{3}_([a-z]*_?)*
var AllValidators = map[string]utils.Validator{
	"001_default_channel": {
		Description: "Ensure defaultChannel is present in list of channels",
		Runner:      validators.ValidateDefaultChannel,
	},
	"002_label_format": {
		Description: "Ensure `label` to follow the format api.openshift.com/addon-<operator-id>",
		Runner:      validators.ValidateAddonLabel,
	},
	"003_csv_present": {
		Description: "Ensure current csv is present in the index image",
		Runner:      validators.ValidateCSVPresent,
	},
}

func getExistingValidatorNames() []string {
	var res []string
	for name := range AllValidators {
		res = append(res, name)
	}
	return res
}

type validatorsFilter struct {
	AllEnabled     bool
	ValidatorNames []string
}

func NewFilter(disabled, enabled string) (*validatorsFilter, error) {
	// empty filter - all validators are enabled
	if disabled == "" && enabled == "" {
		return &validatorsFilter{AllEnabled: true}, nil
	}

	if err := verifyDisabledEnabled(disabled, enabled); err != nil {
		return nil, err
	}

	var validatorNames []string
	if enabled != "" {
		validatorNames = strings.Split(enabled, ",")
	} else {
		validatorNames = getEnabledValidatorNamesFromDisabled(disabled)
	}

	return &validatorsFilter{AllEnabled: false, ValidatorNames: validatorNames}, nil
}

func verifyDisabledEnabled(disabled, enabled string) error {
	// error: mutually exclusive
	if disabled != "" && enabled != "" {
		return errors.New("Can't set both --disabled and --enabled. They are mutually exclusive.")
	}
	var rawNames string
	if enabled != "" {
		rawNames = enabled
	} else {
		rawNames = disabled
	}

	validNames := strings.Join(getExistingValidatorNames(), ",")
	for _, name := range strings.Split(rawNames, ",") {
		if _, ok := AllValidators[name]; !ok {
			return fmt.Errorf("Could not find validator with name %v. Existing validators are %v.", name, validNames)
		}
	}
	return nil
}

func getEnabledValidatorNamesFromDisabled(disabled string) []string {
	var res []string

	allDisabled := make(map[string]bool)
	for _, disabledName := range strings.Split(disabled, ",") {
		allDisabled[disabledName] = true
	}

	for _, name := range getExistingValidatorNames() {
		if _, ok := allDisabled[name]; !ok {
			res = append(res, name)
		}
	}
	return res
}

func (f *validatorsFilter) GetValidators() []utils.Validator {
	var res []utils.Validator
	if f.AllEnabled {
		for _, validator := range AllValidators {
			res = append(res, validator)
		}
	} else {
		for _, name := range f.ValidatorNames {
			// no need to track 'ok' as it was already validated by NewFilter func
			validator := AllValidators[name]
			res = append(res, validator)
		}
	}
	return res
}
