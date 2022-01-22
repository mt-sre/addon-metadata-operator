package extractor

import (
	"context"
	"fmt"
	"sort"

	"github.com/operator-framework/operator-registry/alpha/action"
	"github.com/operator-framework/operator-registry/alpha/model"
	"github.com/sirupsen/logrus"
)

// allBundlesKey - special cacheKey that means "list all bundles for all packages in the indexImage"
const allBundlesKey = "__ALL__"

type DefaultIndexExtractor struct {
	Log   logrus.FieldLogger
	Cache IndexCache
}

// NewIndexExtractor - takes a variadic slice of options to configure an
// index extractor and applies defaults if no appropriate option is given.
func NewIndexExtractor(opts ...IndexExtractorOpt) *DefaultIndexExtractor {
	var extractor DefaultIndexExtractor

	for _, opt := range opts {
		opt(&extractor)
	}

	if extractor.Log == nil {
		extractor.Log = logrus.New()
	}

	if extractor.Cache == nil {
		extractor.Cache = NewIndexMemoryCache()
	}

	extractor.Log = extractor.Log.WithField("source", "indexExtractor")

	return &extractor
}

type IndexExtractorOpt func(e *DefaultIndexExtractor)

func WithIndexCache(cache IndexCache) IndexExtractorOpt {
	return func(e *DefaultIndexExtractor) {
		e.Cache = cache
	}
}

func WithIndexLog(log logrus.FieldLogger) IndexExtractorOpt {
	return func(e *DefaultIndexExtractor) {
		e.Log = log
	}
}

// ExtractBundleImages - returns a sorted list of bundles for a given pkg
func (e *DefaultIndexExtractor) ExtractBundleImages(indexImage string, pkgName string) ([]string, error) {
	e.Log.Debugf("extracting bundles for '%s', matching pkgName '%s'", indexImage, pkgName)
	return e.extractBundleImages(indexImage, pkgName)
}

// ExtractAllBundleImages - returns a sorted list of all bundles for all pkgs
func (e *DefaultIndexExtractor) ExtractAllBundleImages(indexImage string) ([]string, error) {
	e.Log.Debugf("extracting all bundles for '%s'", indexImage)
	return e.extractBundleImages(indexImage, allBundlesKey)
}

// listBundles - return a list of all bundleImages, sorted
func (e *DefaultIndexExtractor) extractBundleImages(indexImage string, cacheKey string) ([]string, error) {
	if bundleImages := e.Cache.GetBundleImages(indexImage, cacheKey); bundleImages != nil {
		e.Log.Debugf("cache hit for '%s'", indexImage)
		return bundleImages, nil
	}

	e.Log.Debugf("cache miss for '%s'", indexImage)
	lb := action.ListBundles{IndexReference: indexImage, PackageName: pkgNameFromCacheKey(cacheKey)}
	data, err := lb.Run(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list bundles with opm: %w", err)
	}

	bundleImages, bundleImagesMap := parseBundles(cacheKey, data.Bundles)
	e.Cache.SetBundleImages(indexImage, bundleImagesMap)

	return bundleImages, nil
}

func pkgNameFromCacheKey(cacheKey string) string {
	if cacheKey == allBundlesKey {
		// a pkg name of "" will list all bundles in an indexImage
		return ""
	}
	return cacheKey
}

func parseBundles(cacheKey string, bundles []model.Bundle) ([]string, map[string][]string) {
	var bundleImages []string
	bundleImagesMap := make(map[string][]string)
	for _, b := range bundles {
		if cacheKey == allBundlesKey || b.Package.Name == cacheKey {
			bundleImagesMap[b.Package.Name] = append(bundleImagesMap[b.Package.Name], b.Image)
		}
		bundleImages = append(bundleImages, b.Image)
	}
	sort.Strings(bundleImages)
	return bundleImages, bundleImagesMap
}
