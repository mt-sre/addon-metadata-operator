package extractor

import (
	"sync"

	"github.com/operator-framework/operator-registry/pkg/registry"
)

type bundleMemoryCache struct {
	encoder BundleEncoder

	lock  sync.RWMutex
	store map[string][]byte
}

// NewBundleMemoryCache - in-memory cache for bundle, using an encoder
func NewBundleMemoryCache(encoder BundleEncoder) BundleCache {
	return &bundleMemoryCache{
		encoder: encoder,
		store:   make(map[string][]byte),
	}
}

func (c *bundleMemoryCache) Get(bundleImage string) (*registry.Bundle, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if data, ok := c.store[bundleImage]; ok {
		bundle, err := c.encoder.Decode(data)
		if err != nil {
			return nil, err
		}

		return hackSetCacheStaleToTrue(bundle), nil
	}

	return nil, nil
}

// hack to set b.cacheStale to true otherwise we can't access the csv of the
// underlying bundle. This is a bug on OPM's side, which does not support
// serialization/deserialization of their bundles.
// https://github.com/operator-framework/operator-registry/blob/master/pkg/registry/bundle.go#L103
func hackSetCacheStaleToTrue(b *registry.Bundle) *registry.Bundle {
	if len(b.Objects) == 0 {
		return b
	}

	obj := b.Objects[0]
	b.Objects = b.Objects[1:]
	b.Add(obj)

	return b
}

func (c *bundleMemoryCache) Set(bundleImage string, bundle *registry.Bundle) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// bundles are supposed to be immutable, so we can save a few CPU cycles by
	// avoiding re-encoding bundles and unnecessarily overwriting the cache
	if _, ok := c.store[bundleImage]; ok {
		return nil
	}

	data, err := c.encoder.Encode(bundle)
	if err != nil {
		return err
	}

	c.store[bundleImage] = data

	return nil
}
