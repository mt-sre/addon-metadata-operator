package utils

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

type Validator struct {
	Name        string
	Description string
	Runner      ValidateFunc
}

type ValidatorTest interface {
	Name() string
	Run(MetaBundle) (bool, string, error)
	SucceedingCandidates() []MetaBundle
	FailingCandidates() []MetaBundle
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

// ValidateFunc - returns a triple consisting of:
// 1. bool       - true if MetaBundle validation was successful
// 2. failureMsg - "" if the result was true, else information about why the validation failed
// 3. error      - report any error that happened in the validation code
type ValidateFunc func(mb MetaBundle) (bool, string, error)

type MetaBundle struct {
	AddonMeta *v1alpha1.AddonMetadataSpec
	Bundles   []registry.Bundle
}

func NewMetaBundle(addonMeta *v1alpha1.AddonMetadataSpec, bundles []registry.Bundle) *MetaBundle {
	return &MetaBundle{
		AddonMeta: addonMeta,
		Bundles:   bundles,
	}
}
