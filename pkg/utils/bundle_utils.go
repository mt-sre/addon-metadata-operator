package utils

import (
	"errors"
	"io/ioutil"

	"github.com/operator-framework/operator-registry/pkg/registry"
	log "github.com/sirupsen/logrus"
)

const AllAddonsIdentifier = "all"

var (
	bundleParser BundleParser
)

func init() {
	bundleParser = DefaultBundleParser{}
}

func ExtractAndParseAddons(indexImage, addonIdentifier string) ([]registry.Bundle, error) {
	if indexImage == "" {
		return []registry.Bundle{}, errors.New("Missing index image!")
	}

	indexImageExtractor := DefaultIndexImageExtractor{
		downloadPath: DefaultDownloadPath,
		cacheDir:     DefaultCacheDir,
		indexImage:   indexImage,
	}

	key := indexImageExtractor.CacheKey(indexImage, addonIdentifier)
	if err := extractAddons(indexImageExtractor, key, indexImage); err != nil {
		return []registry.Bundle{}, err
	}
	if addonIdentifier == AllAddonsIdentifier {
		// Parse all adddons in the extraction path
		return parseAllAddons(indexImageExtractor, indexImageExtractor.ExtractionPath())
	}
	// Parse only a specific addon in the extraction path
	manifestsDir := indexImageExtractor.ManifestsPath(addonIdentifier)
	return bundleParser.ParseBundles(addonIdentifier, manifestsDir)
}

func extractAddons(indexImageExtractor DefaultIndexImageExtractor, cacheKey, indexImage string) error {
	if !indexImageExtractor.CacheHit(cacheKey) {
		if err := indexImageExtractor.ExtractBundlesFromImage(indexImage, indexImageExtractor.ExtractionPath()); err != nil {
			return err
		}
		if err := indexImageExtractor.WriteToCache(cacheKey); err != nil {
			log.Warnln("Failed to cache the index image bundles")
		}
	}
	return nil
}

func parseAllAddons(indexImageExtractor DefaultIndexImageExtractor, addonsDir string) ([]registry.Bundle, error) {
	operatorsFound, err := allOperatorsFound(addonsDir)
	if err != nil {
		return []registry.Bundle{}, err
	}

	var allBundles []registry.Bundle
	for _, operator := range operatorsFound {
		manifestsDir := indexImageExtractor.ManifestsPath(operator)
		bundles, err := bundleParser.ParseBundles(operator, manifestsDir)
		if err != nil {
			return []registry.Bundle{}, err
		}
		allBundles = append(allBundles, bundles...)
	}

	return allBundles, nil
}

func allOperatorsFound(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, err
	}
	operatorsFound := []string{}
	for _, file := range files {
		if file.IsDir() {
			operatorsFound = append(operatorsFound, file.Name())
		}
	}
	return operatorsFound, nil
}
