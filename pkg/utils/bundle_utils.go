package utils

import (
	"errors"
	"io/ioutil"

	"github.com/operator-framework/operator-registry/pkg/registry"
	log "github.com/sirupsen/logrus"
)

const allAddonsKey = "all"

var (
	bundleParser BundleParser
)

func init() {
	bundleParser = DefaultBundleParser{}
}

func ExtractAndParseAddon(indexImage, addonName string) ([]registry.Bundle, error) {
	if indexImage == "" {
		return []registry.Bundle{}, errors.New("Missing index image!")
	}

	indexImageExtractor := DefaultIndexImageExtractor{
		downloadPath: defaultDownloadPath,
		cacheDir:     defaultCacheDir,
		indexImage:   indexImage,
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

func ExtractAndParseAllAddons(indexImage string) ([]registry.Bundle, error) {
	if indexImage == "" {
		return []registry.Bundle{}, errors.New("Missing index image!")
	}
	indexImageExtractor := DefaultIndexImageExtractor{
		downloadPath: defaultDownloadPath,
		cacheDir:     defaultCacheDir,
		indexImage:   indexImage,
	}
	key := indexImageExtractor.CacheKey(indexImage, allAddonsKey)

	if !indexImageExtractor.CacheHit(key) {
		if err := indexImageExtractor.ExtractBundlesFromImage(indexImage, indexImageExtractor.ExtractionPath()); err != nil {
			return []registry.Bundle{}, err
		}
		if err := indexImageExtractor.WriteToCache(key); err != nil {
			log.Warnln("Failed to cache the index image bundles")
		}
	}
	operatorsFound, err := allOperatorsFound(indexImageExtractor.ExtractionPath())
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
