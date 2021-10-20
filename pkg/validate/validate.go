package validate

import (
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

func Validate(mb *utils.MetaBundle, runMeta bool) []error {
	errs := []error{}

	if runMeta {
		validators := GetAllMetaValidators()

		printMetaHeading()

		for _, validator := range validators {
			fmt.Printf("\r%s\t\t", validator.Description)
			success, err := validator.Runner(mb)
			if err != nil {
				errs = append(errs, err)
				printErrorMessage(validator.Description)
			} else if !success {
				printFailureMessage(validator.Description)
			} else {
				printSuccessMessage(validator.Description)
			}
			fmt.Println()
		}
	}
	return errs
}
