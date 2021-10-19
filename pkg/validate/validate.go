package validate

import (
	"github.com/go-playground/validator"
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
)

// use a single instance of MetadataValidator as per docs because it
// caches struct info
var metaBundlevalidator *validator.Validate

type MetaBundle struct {
	AddonMeta *v1alpha1.AddonMetadataSpec
	// TODO: add field for correspinding bundle
}

// TODO: This will return a MetaBundle with corresponding bundle
func NewMetaBundle(addonMeta *v1alpha1.AddonMetadataSpec) *MetaBundle {
	return &MetaBundle{
		AddonMeta: addonMeta,
	}
}

func getMetaBundleValidator(registerMeta bool) *validator.Validate {
	if metaBundlevalidator != nil {
		return metaBundlevalidator
	}

	if registerMeta {
		bundleValidations := GetAllMetaValidators()
		for _, validation := range bundleValidations {
			metaBundlevalidator.RegisterStructValidation(validation.Runner, MetaBundle{})
		}
	}
	return metaBundlevalidator
}

func (mb *MetaBundle) Validate(runMeta bool) *[]validator.FieldError {
	validate := getMetaBundleValidator(runMeta)

	if err := validate.Struct(mb); err != nil {
		res := []validator.FieldError{}
		for _, fieldError := range err.(validator.ValidationErrors) {
			res = append(res, fieldError)
		}
		return &res
	}
	return nil
}
