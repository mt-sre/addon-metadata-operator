package validate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"

	"github.com/alexeyco/simpletable"
)

// ValidateCLI - run all validators on a metaBundle struct
func ValidateCLI(mb types.MetaBundle, filter *validatorsFilter) (bool, []error) {
	errs := []error{}
	allSuccess := true

	t := simpletable.New()
	t.Header = getTableHeaders()

	for _, v := range filter.GetValidators() {
		row := newSuccessTableRow(v)
		res := v.Runner(mb)
		if res.IsError() {
			errs = append(errs, res.Error)
			row = newErrorTableRow(v, res.Error)
		} else if !res.IsSuccess() {
			allSuccess = false
			row = newFailedTableRow(v, res.FailureMsg)
		}
		t.Body.Cells = append(t.Body.Cells, row)
	}
	t.SetStyle(simpletable.StyleCompactLite)
	fmt.Printf("%v\n\n", t.String())
	fmt.Println("Please consult corresponding validator wikis: https://github.com/mt-sre/addon-metadata-operator/wiki/<code>.")
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

func (f *validatorsFilter) GetValidators() []types.Validator {
	var res []types.Validator
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
