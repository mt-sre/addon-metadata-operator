package validate

import (
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
)

type Validator struct {
	Description string
	Runner      ValidateFunc
}

type ValidateFunc func(mb *MetaBundle) error

type MetaBundle struct {
	AddonMeta *v1alpha1.AddonMetadataSpec
	// TODO: add field for corresponding bundle
}

// TODO: This will return a MetaBundle with corresponding bundle
func NewMetaBundle(addonMeta *v1alpha1.AddonMetadataSpec) *MetaBundle {
	return &MetaBundle{
		AddonMeta: addonMeta,
	}
}

func (mb *MetaBundle) Validate(runMeta bool) []error {
	errs := []error{}

	if runMeta {
		validators := GetAllMetaValidators()

		printMetaHeading()

		for _, validator := range validators {
			fmt.Printf("\r%s\t\t", validator.Description)
			err := validator.Runner(mb)
			if err != nil {
				errs = append(errs, err)
				printFailureMessage(validator.Description)
			} else {
				printSuccessMessage(validator.Description)
			}
			fmt.Println()
		}
	}
	return errs
}
