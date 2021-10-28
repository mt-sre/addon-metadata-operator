package utils

import (
	"fmt"
)

// PrintValidationErrors - helper to pretty print validationErrors
func PrintValidationErrors(errs []error) {
	fmt.Printf("\n%s\n", Red("Failed with the following errors:"))
	for _, err := range errs {
		fmt.Printf("%s\n", err.Error())
	}
}
