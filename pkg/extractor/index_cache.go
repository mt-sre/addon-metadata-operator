package extractor

import (
	"sync"
)

type indexMemoryCache struct {
	sync.RWMutex
	store map[string]map[string][]string
}

// NewIndexMemoryCache - in-memory cache for bundleImages. Indexes first by indexImage,
// then by pkgName.
func NewIndexMemoryCache() IndexCache {
	return &indexMemoryCache{store: make(map[string]map[string][]string)}
}

func (c *indexMemoryCache) GetBundleImages(indexImage string, cacheKey string) []string {
	c.RLock()
	defer c.RUnlock()

	var res []string
	if data, ok := c.store[indexImage]; ok {
		for pkgName, bundles := range data {
			if cacheKey == allBundlesKey || pkgName == cacheKey {
				res = append(res, bundles...)
			}
		}
	}
	if len(res) == 0 {
		return nil
	}
	return res
}

// SetBundleImages - indexImages are mutable, so we always need to update the cache
func (c *indexMemoryCache) SetBundleImages(indexImage string, pkgBundlesMap map[string][]string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.store[indexImage]; ok {
		for pkgName, bundles := range pkgBundlesMap {
			c.store[indexImage][pkgName] = bundles
		}
	} else {
		c.store[indexImage] = pkgBundlesMap
	}
}
