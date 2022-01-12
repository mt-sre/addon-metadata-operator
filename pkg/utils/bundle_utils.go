package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/operator-framework/operator-registry/pkg/registry"
	log "github.com/sirupsen/logrus"
)

const AllAddonsIdentifier = "all"

var (
	bundleParser types.BundleParser
	lock         = sync.Mutex{}
)

func init() {
	bundleParser = DefaultBundleParser{}
}

func ExtractAndParseAddons(indexImage, addonIdentifier string) ([]registry.Bundle, error) {
	// TODO (sblaisdo) - when migrating to server code, need something more complex as this
	// is a global lock that will prevent goroutines from extracting indexImage in parallel
	// avoid race condition when extracting same indexImage in two goroutines
	lock.Lock()
	defer lock.Unlock()

	if indexImage == "" {
		return []registry.Bundle{}, errors.New("Missing index image!")
	}

	// ignore tagless images used by addons in the process of on-boarding
	if parts := strings.SplitN(indexImage, ":", 2); len(parts) < 2 {
		return []registry.Bundle{}, nil
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

	if len(allBundles) == 0 {
		return []registry.Bundle{}, errors.New("No bundles were found.")
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

// GetBundleNameVersion - useful for validation error reporting
func GetBundleNameVersion(b registry.Bundle) (string, error) {
	version, err := b.Version()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v:%v", b.Name, version), nil
}
