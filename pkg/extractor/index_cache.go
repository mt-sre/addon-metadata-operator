package extractor

import (
	"sort"
	"sync"
)

type indexMemoryCache struct {
	sync.RWMutex
	store map[string]map[string][]string
}

// Simple in-memory cache to avoid extracting indexImages twice
// Indexes first by indexImage, then by packageName
func NewIndexMemoryCache() IndexExtractorCache {
	return &indexMemoryCache{store: make(map[string]map[string][]string)}
}

func (c *indexMemoryCache) GetBundles(indexImage string, cacheKey string) []string {
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
	sort.Strings(res)
	return res
}

func (c *indexMemoryCache) SetBundles(indexImage string, pkgBundlesMap map[string][]string) {
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
