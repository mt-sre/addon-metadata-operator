package utils

import (
	"errors"

	"github.com/operator-framework/operator-registry/pkg/registry"
	log "github.com/sirupsen/logrus"
)

var (
	bundleParser        BundleParser
	indexImageExtractor IndexImageExtractor
)

func init() {
	indexImageExtractor = DefaultIndexImageExtractor{
		downloadPath:  defaultDownloadPath,
		cacheLocation: defaultCacheLocation,
	}
	bundleParser = DefaultBundleParser{}

}

func ExtractAndParse(indexImage, addonName string) ([]registry.Bundle, error) {
	if indexImage == "" {
		return []registry.Bundle{}, errors.New("Missing index image!")
	}
	key := indexImageExtractor.CacheKey(indexImage, addonName)
	if !indexImageExtractor.CacheHit(key) {
		if err := indexImageExtractor.ExtractBundlesFromImage(indexImage, indexImageExtractor.ExtractionPath()); err != nil {
			return []registry.Bundle{}, err
		}
		if err := indexImageExtractor.WriteToCache(key); err != nil {
			log.Warnln("Failed to cache the index image bundles")
		}
	}
	manifestsDir := indexImageExtractor.ManifestsPath(addonName)
	return bundleParser.ParseBundles(addonName, manifestsDir)
}
