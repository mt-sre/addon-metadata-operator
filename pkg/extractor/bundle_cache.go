package extractor

import (
	"errors"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/operator"
)

// NewBundleCacheImpl returns an initialized BundleCacheImpl. A
// variadic slice of options can be used to alter the behavior
// of the cache.
func NewBundleCacheImpl(opts ...BundleCacheImplOption) *BundleCacheImpl {
	var cfg BundleCacheImplConfig

	cfg.Option(opts...)
	cfg.Default()

	return &BundleCacheImpl{
		cfg: cfg,
	}
}

type BundleCacheImpl struct {
	cfg BundleCacheImplConfig
}

var ErrInvalidBundleData = errors.New("invalid bundle data")

func (c *BundleCacheImpl) GetBundle(img string) (*operator.Bundle, error) {
	data, ok := c.cfg.Store.Read(img)
	if !ok {
		return nil, nil
	}

	bundle, ok := data.(operator.Bundle)
	if !ok {
		return nil, ErrInvalidBundleData
	}

	return &bundle, nil
}

func (c *BundleCacheImpl) SetBundle(img string, bundle operator.Bundle) error {
	if err := c.cfg.Store.Write(img, bundle); err != nil {
		return fmt.Errorf("writing bundle data: %w", err)
	}

	return nil
}

type BundleCacheImplConfig struct {
	Store Store
}

func (c *BundleCacheImplConfig) Option(opts ...BundleCacheImplOption) {
	for _, opt := range opts {
		opt.ConfigureBundleCacheImpl(c)
	}
}

func (c *BundleCacheImplConfig) Default() {
	if c.Store == nil {
		c.Store = NewThreadSafeStore()
	}
}

type BundleCacheImplOption interface {
	ConfigureBundleCacheImpl(*BundleCacheImplConfig)
}
