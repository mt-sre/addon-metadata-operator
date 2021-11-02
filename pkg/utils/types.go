package utils

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

type Validator struct {
	Description string
	Runner      ValidateFunc
}

type IndexImageExtractor interface {
	ExtractBundlesFromImage(indexImage string, extractTo string) error
	CacheKey(indexImage, addonName string) string
	CacheHit(key string) bool
	ExtractionPath() string
	ManifestsPath(addonName string) string
	CacheLocation() string
	WriteToCache(value string) error
}

type BundleParser interface {
	ParseBundles(addonName string, manifestsPath string) ([]registry.Bundle, error)
}

type ValidateFunc func(mb *MetaBundle) (bool, error)

type MetaBundle struct {
	AddonMeta *v1alpha1.AddonMetadataSpec
	Bundles   []registry.Bundle
}

// TODO: This will return a MetaBundle with corresponding bundle
func NewMetaBundle(addonMeta *v1alpha1.AddonMetadataSpec, bundles []registry.Bundle) *MetaBundle {
	return &MetaBundle{
		AddonMeta: addonMeta,
		Bundles:   bundles,
	}
}
