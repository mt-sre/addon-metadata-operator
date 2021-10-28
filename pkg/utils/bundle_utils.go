package utils

import (
	"errors"

	"github.com/operator-framework/operator-registry/pkg/registry"
	log "github.com/sirupsen/logrus"
)

var bundleParser BundleParser
var indexImageExtractor IndexImageExtractor

func init() {
	indexImageExtractor = DefaultIndexImageExtractor{
		downloadPath:  DefaultDownloadPath,
		cacheLocation: DefaultCacheLocation,
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
	}
	if err := indexImageExtractor.WriteToCache(key); err != nil {
		log.Warnln("Failed to cache the index image bundles")
	}
	manifestsDir := indexImageExtractor.ManifestsPath(addonName)
	result, err := bundleParser.ParseBundles(addonName, manifestsDir)
	if err != nil {
		return []registry.Bundle{}, err
	}
	return result, nil
}
