package types

import "github.com/operator-framework/operator-registry/pkg/registry"

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
