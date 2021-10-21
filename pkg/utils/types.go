package utils

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
)

type Validator struct {
	Description string
	Runner      ValidateFunc
}

type ValidateFunc func(mb *MetaBundle) (bool, error)

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
