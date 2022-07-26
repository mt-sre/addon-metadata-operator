package extractor

import (
	"context"
	"fmt"
	"sort"

	"github.com/operator-framework/operator-registry/alpha/action"
	"github.com/operator-framework/operator-registry/alpha/model"
	"github.com/sirupsen/logrus"
)

// IndexCache provides a cache of index images which store related bundles
// by package name.
type IndexCache interface {
	// GetBundleImages retrieves bundles image names for a particular
	// indexImage and package combination.
	GetBundleImages(indexImage string, pkgName string) ([]string, error)
	// SetBundleImages stores a map from package name to bundle images
	// for a particular indexImage. An error is returned if the data
	// cannot be written.
	SetBundleImages(indexImage string, bundleImagesMap map[string][]string) error
}

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
		extractor.Cache = NewIndexCacheImpl()
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

// listBundles - return a list of all bundleImages. Need to sort bundleImages
// everytime as order might not be preserved in the cache.
func (e *DefaultIndexExtractor) extractBundleImages(indexImage string, cacheKey string) ([]string, error) {
	bundleImages, err := e.Cache.GetBundleImages(indexImage, cacheKey)
	if err != nil {
		e.Log.Warnf("getting bundle images from cache: %w", err)
	}

	if bundleImages != nil {
		e.Log.Debugf("cache hit for '%s'", indexImage)
		return sortedBundleImages(bundleImages), nil
	}

	e.Log.Debugf("cache miss for '%s'", indexImage)
	lb := action.ListBundles{IndexReference: indexImage, PackageName: pkgNameFromCacheKey(cacheKey)}
	data, err := lb.Run(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list bundles with opm: %w", err)
	}

	bundleImages, bundleImagesMap := parseBundles(cacheKey, data.Bundles)

	if err := e.Cache.SetBundleImages(indexImage, bundleImagesMap); err != nil {
		e.Log.Warnf("caching bundle images: %w", err)
	}

	return sortedBundleImages(bundleImages), nil
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
	return bundleImages, bundleImagesMap
}

func sortedBundleImages(bundleImages []string) []string {
	sort.Strings(bundleImages)
	return bundleImages
}
