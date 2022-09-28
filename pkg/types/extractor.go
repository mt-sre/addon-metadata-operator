package types

type IndexImageExtractor interface {
	ExtractBundlesFromImage(indexImage string, extractTo string) error
	CacheKey(indexImage, addonName string) string
	CacheHit(key string) bool
	ExtractionPath() string
	ManifestsPath(addonName string) string
	CacheLocation() string
	WriteToCache(value string) error
}
