package utils

import (
	"fmt"

	"github.com/go-playground/validator"
)

// PrintValidationErrors - helper to pretty print validationErrors
func PrintValidationErrors(fieldErrors *[]validator.FieldError) {
	for _, err := range *fieldErrors {
		fmt.Println(err)
	}
}
