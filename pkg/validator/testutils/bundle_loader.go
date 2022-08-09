package testutils

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/stretchr/testify/require"
)

func NewBundlerLoader(t *testing.T) *BundleLoader {
	return &BundleLoader{
		t: t,
	}
}

type BundleLoader struct {
	t *testing.T
}

func (l *BundleLoader) LoadFromCSV(path string, opts ...LoadOption) *registry.Bundle {
	l.t.Helper()

	var cfg loadConfig

	cfg.Option(opts...)
	cfg.Default()

	bundle, err := testutils.NewBundle(cfg.BundleName, path)
	require.NoError(l.t, err)

	if cfg.PackageName != "" {
		bundle.Annotations.PackageName = cfg.PackageName
	}

	return bundle
}

type loadConfig struct {
	BundleName  string
	PackageName string
}

func (c *loadConfig) Option(opts ...LoadOption) {
	for _, opt := range opts {
		opt.ConfigureBundleLoad(c)
	}
}

func (c *loadConfig) Default() {
	if c.BundleName == "" {
		c.BundleName = "test-bundle"
	}
}

type LoadOption interface {
	ConfigureBundleLoad(*loadConfig)
}

type WithBundleName string

func (w WithBundleName) ConfigureBundleLoad(c *loadConfig) {
	c.BundleName = string(w)
}

type WithPackageName string

func (w WithPackageName) ConfigureBundleLoad(c *loadConfig) {
	c.PackageName = string(w)
}
