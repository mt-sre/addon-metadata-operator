package testutils

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/operator"
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

func (l *BundleLoader) LoadFromCSV(path string, opts ...LoadOption) operator.Bundle {
	l.t.Helper()

	var cfg loadConfig

	cfg.Option(opts...)
	cfg.Default()

	regBundle, err := testutils.NewBundle(cfg.BundleName, path)
	require.NoError(l.t, err)

	if cfg.PackageName != "" {
		regBundle.Annotations.PackageName = cfg.PackageName
	}

	bundle, err := operator.NewBundleFromRegistryBundle(*regBundle)
	require.NoError(l.t, err)

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
