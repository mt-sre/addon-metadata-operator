package validate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/validators"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

// Validate - run all validators on a metaBundle struct
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

func getExistingValidatorCodes() []string {
	var res []string
	for code := range validators.Registry.All() {
		res = append(res, code)
	}
	return res
}

type validatorsFilter struct {
	AllEnabled     bool
	ValidatorCodes []string
}

func NewFilter(disabled, enabled string) (*validatorsFilter, error) {
	// empty filter - all validators are enabled
	if disabled == "" && enabled == "" {
		return &validatorsFilter{AllEnabled: true}, nil
	}

	if err := verifyDisabledEnabled(disabled, enabled); err != nil {
		return nil, err
	}

	var validatorCodes []string
	if enabled != "" {
		validatorCodes = strings.Split(enabled, ",")
	} else {
		validatorCodes = getEnabledValidatorCodesFromDisabled(disabled)
	}

	return &validatorsFilter{AllEnabled: false, ValidatorCodes: validatorCodes}, nil
}

func verifyDisabledEnabled(disabled, enabled string) error {
	// error: mutually exclusive
	if disabled != "" && enabled != "" {
		return errors.New("Can't set both --disabled and --enabled. They are mutually exclusive.")
	}
	var rawCodes string
	if enabled != "" {
		rawCodes = enabled
	} else {
		rawCodes = disabled
	}

	validCodes := strings.Join(getExistingValidatorCodes(), ",")
	for _, code := range strings.Split(rawCodes, ",") {
		if _, ok := validators.Registry.Get(code); !ok {
			return fmt.Errorf("Could not find validator with code %v. Existing validators are %v.", code, validCodes)
		}
	}
	return nil
}

func getEnabledValidatorCodesFromDisabled(disabled string) []string {
	var res []string

	allDisabled := make(map[string]bool)
	for _, disabledCode := range strings.Split(disabled, ",") {
		allDisabled[disabledCode] = true
	}

	for _, code := range getExistingValidatorCodes() {
		if _, ok := allDisabled[code]; !ok {
			res = append(res, code)
		}
	}
	return res
}

func (f *validatorsFilter) GetValidators() []utils.Validator {
	var res []utils.Validator
	if f.AllEnabled {
		for _, validator := range validators.Registry.All() {
			res = append(res, validator)
		}
	} else {
		for _, code := range f.ValidatorCodes {
			// no need to track 'ok' as it was already validated by NewFilter func
			validator, _ := validators.Registry.Get(code)
			res = append(res, validator)
		}
	}
	return res
}
