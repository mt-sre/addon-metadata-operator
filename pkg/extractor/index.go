package extractor

import (
	"context"
	"sort"

	"github.com/operator-framework/operator-registry/alpha/action"
	"github.com/operator-framework/operator-registry/alpha/model"
)

// allBundlesKey - special cacheKey that means "list all bundles for all packages in the indexImage"
const allBundlesKey = "__ALL__"

type indexExtractor struct {
	IndexExtractorCache
}

func NewIndexExtractor(cache IndexExtractorCache) IndexExtractor {
	return &indexExtractor{cache}
}

// ListBundlesFromPackage - returns a sorted list of bundles for a given pkg
func (i *indexExtractor) ListBundlesFromPackage(indexImage string, pkgName string) ([]string, error) {
	return i.listBundlesWithCache(indexImage, pkgName)
}

// ListAllBundles - returns a sorted list of all bundles for all pkgs
func (i *indexExtractor) ListAllBundles(indexImage string) ([]string, error) {
	return i.listBundlesWithCache(indexImage, allBundlesKey)
}

// listBundles - return a list of all bundle images, sorted
func (i *indexExtractor) listBundlesWithCache(indexImage string, cacheKey string) ([]string, error) {
	if bundles := i.GetBundles(indexImage, cacheKey); bundles != nil {
		return bundles, nil
	}
	pkgName := pkgNameFromCacheKey(cacheKey)
	lb := action.ListBundles{IndexReference: indexImage, PackageName: pkgName}
	data, err := lb.Run(context.Background())
	if err != nil {
		return nil, err
	}
	return i.extractAndCache(indexImage, cacheKey, data.Bundles), nil
}

func pkgNameFromCacheKey(cacheKey string) string {
	if cacheKey == allBundlesKey {
		// a pkg name of "" will list all bundles in an indexImage
		return ""
	}
	return cacheKey
}

func (i *indexExtractor) extractAndCache(indexImage string, cacheKey string, bundles []model.Bundle) []string {
	var res []string
	pkgBundlesMap := make(map[string][]string)
	for _, b := range bundles {
		if cacheKey == allBundlesKey || b.Package.Name == cacheKey {
			pkgBundlesMap[b.Package.Name] = append(pkgBundlesMap[b.Package.Name], b.Image)
		}
		res = append(res, b.Image)
	}
	i.SetBundles(indexImage, pkgBundlesMap)
	sort.Strings(res)
	return res
}
