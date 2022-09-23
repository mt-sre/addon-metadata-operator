package extractor

import (
	"errors"
	"fmt"
)

type IndexCacheImpl struct {
	cfg IndexCacheImplConfig
}

// NewIndexCacheImpl returns an initialized IndexImage cache which stores
// the package names and related bundle images for an IndexImage.
// A variadic slice of options may be provided to alter the cache behavior.
func NewIndexCacheImpl(opts ...IndexCacheImplOption) *IndexCacheImpl {
	var cfg IndexCacheImplConfig

	cfg.Option(opts...)
	cfg.Default()

	return &IndexCacheImpl{
		cfg: cfg,
	}
}

var ErrInvalidIndexData = errors.New("invalid index data")

func (c *IndexCacheImpl) GetBundleImages(indexImage string, cacheKey string) ([]string, error) {
	data, ok := c.cfg.Store.Read(indexImage)
	if !ok {
		return nil, nil
	}

	pkgs, ok := data.(map[string][]string)
	if !ok {
		return nil, ErrInvalidIndexData
	}

	var res []string

	for pkgName, bundles := range pkgs {
		if cacheKey != allBundlesKey && pkgName != cacheKey {
			continue
		}

		res = append(res, bundles...)
	}

	return res, nil
}

func (c *IndexCacheImpl) SetBundleImages(indexImage string, pkgBundlesMap map[string][]string) error {
	data := make(map[string][]string, len(pkgBundlesMap))

	for pkg, bundles := range pkgBundlesMap {
		data[pkg] = bundles
	}

	if err := c.cfg.Store.Write(indexImage, data); err != nil {
		return fmt.Errorf("writing data: %w", err)
	}

	return nil
}

type IndexCacheImplConfig struct {
	Store Store
}

func (c *IndexCacheImplConfig) Option(opts ...IndexCacheImplOption) {
	for _, opt := range opts {
		opt.ConfigureIndexCacheImpl(c)
	}
}

func (c *IndexCacheImplConfig) Default() {
	if c.Store == nil {
		c.Store = NewThreadSafeStore()
	}
}

type IndexCacheImplOption interface {
	ConfigureIndexCacheImpl(*IndexCacheImplConfig)
}
