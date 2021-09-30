package v1alpha1

import (
	"github.com/go-playground/validator"
)

// use a single instance of MetadataValidator as per docs because it
// caches struct info
var metadataValidator *validator.Validate

// getMetadataValidator - returns the MetadataValidator Singleton
func getMetadataValidator() *validator.Validate {
	if metadataValidator != nil {
		return metadataValidator
	}

	// instantiate the singleton
	metadataValidator = validator.New()

	// register validation functions here
	// TODO
	// ...

	return metadataValidator
}

// Validate - abstracts details behind AddonMetadata validation
// returns pointer to list of FieldError so we can nil the result
func (a *AddonMetadata) Validate() *[]validator.FieldError {
	metadataValidator = getMetadataValidator()

	if err := metadataValidator.Struct(a); err != nil {
		res := []validator.FieldError{}
		for _, fieldError := range err.(validator.ValidationErrors) {
			res = append(res, fieldError)
		}
		return &res
	}

	// no errors
	return nil
}
