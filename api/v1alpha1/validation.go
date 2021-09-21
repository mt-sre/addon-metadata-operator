package v1alpha1

import "github.com/go-playground/validator"

// use a single instance of MetadataValidator, it caches struct info
var metadataValidator *validator.Validate

// GetMetadataValidator - returns the MetadataValidator Singleton
func GetMetadataValidator() *validator.Validate {
	if metadataValidator != nil {
		return metadataValidator
	}

	// instantiate the singleton
	metadataValidator = validator.New()
	return metadataValidator
}
