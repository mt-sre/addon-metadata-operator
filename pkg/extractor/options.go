package extractor

// WithStore provides the given Store implementation
// to cache implementation for object storage.
type WithStore struct{ Store }

func (w WithStore) ConfigureIndexCacheImpl(c *IndexCacheImplConfig) {
	c.Store = w.Store
}

func (w WithStore) ConfigureBundleCacheImpl(c *BundleCacheImplConfig) {
	c.Store = w.Store
}
