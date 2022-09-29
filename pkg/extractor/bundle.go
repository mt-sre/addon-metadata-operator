package extractor

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mt-sre/addon-metadata-operator/pkg/operator"
	"github.com/operator-framework/operator-registry/pkg/image"
	"github.com/operator-framework/operator-registry/pkg/image/containerdregistry"
	opmbundle "github.com/operator-framework/operator-registry/pkg/lib/bundle"
	"github.com/sirupsen/logrus"
)

// BundleCache provides a cache of OPM bundles which are referenced by
// bundleImage name.
type BundleCache interface {
	// GetBundle returns a bundle for the given image. An error is
	// returned if the bundle cannot be retrieved or the data is
	// corrupted.
	GetBundle(img string) (*operator.Bundle, error)
	// SetBundle caches a bundle for the given image. An error
	// is returned if the bundle cannot be cached.
	SetBundle(img string, bundle operator.Bundle) error
}

type DefaultBundleExtractor struct {
	Log     logrus.FieldLogger
	Cache   BundleCache
	Timeout time.Duration
}

func NewBundleExtractor(opts ...BundleExtractorOpt) *DefaultBundleExtractor {
	const defaultTimeout = 60 * time.Second

	extractor := DefaultBundleExtractor{
		Timeout: defaultTimeout,
	}
	for _, opt := range opts {
		opt(&extractor)
	}

	if extractor.Log == nil {
		extractor.Log = logrus.New()
	}

	if extractor.Cache == nil {
		extractor.Cache = NewBundleCacheImpl()
	}

	extractor.Log = extractor.Log.WithField("source", "bundleExtractor")
	return &extractor
}

type BundleExtractorOpt func(e *DefaultBundleExtractor)

func WithBundleCache(cache BundleCache) BundleExtractorOpt {
	return func(e *DefaultBundleExtractor) {
		e.Cache = cache
	}
}

func WithBundleLog(log logrus.FieldLogger) BundleExtractorOpt {
	return func(e *DefaultBundleExtractor) {
		e.Log = log
	}
}

func WithBundleTimeout(timeout time.Duration) BundleExtractorOpt {
	return func(e *DefaultBundleExtractor) {
		e.Timeout = timeout
	}
}

func (e *DefaultBundleExtractor) Extract(ctx context.Context, bundleImage string) (operator.Bundle, error) {
	cachedBundle, err := e.Cache.GetBundle(bundleImage)
	if err != nil {
		e.Log.Warnf("retrieving bundle %q from cache: %w", bundleImage, err)
	}

	if cachedBundle != nil {
		e.Log.Debugf("cache hit for %q", bundleImage)
		return *cachedBundle, nil
	}

	e.Log.Debugf("cache miss for '%s'", bundleImage)
	tmpDirs, err := createTempDirs()
	if err != nil {
		return operator.Bundle{}, err
	}
	defer func() {
		if err := tmpDirs.CleanUp(); err != nil {
			e.Log.Errorf("cleaning up tmpDirs: %w", err)
		}
	}()

	if err := e.unpackAndValidateBundle(ctx, bundleImage, tmpDirs); err != nil {
		return operator.Bundle{}, fmt.Errorf("unpacking and validating bundle: %w", err)
	}

	bundle, err := operator.NewBundleFromDirectory(tmpDirs["bundle"])
	if err != nil {
		return operator.Bundle{}, err
	}

	bundle.BundleImage = bundleImage // not set by OPM

	if err := e.Cache.SetBundle(bundleImage, bundle); err != nil {
		e.Log.Warnf("caching bundle %q: %w", bundleImage, err)
	}

	return bundle, nil
}

// unpackAndValidateBundle - Unpacks the content of an operator bundle into a temp directory
// and validates the extracted bundle.
// Reference: https://github.com/operator-framework/operator-registry/blob/master/cmd/opm/alpha/bundle/unpack.go
func (e *DefaultBundleExtractor) unpackAndValidateBundle(ctx context.Context, bundleImage string, tmpDirs tempDirs) error {
	e.Log.Debugf("unpacking bundleImage '%s' to '%s'", bundleImage, tmpDirs["bundle"])

	ctx, cancel := context.WithTimeout(ctx, e.Timeout)
	defer cancel()

	registry, err := containerdregistry.NewRegistry(
		containerdregistry.SkipTLSVerify(false),
		containerdregistry.WithLog(e.Log.(*logrus.Entry)),
		// need a new cache dir for each registry to avoid data races and
		// having the default "cache/ingest" dir removed from under our feet
		containerdregistry.WithCacheDir(tmpDirs["containerd"]),
	)
	if err != nil {
		return err
	}
	defer func() {
		// ensure cleanup of registry resources, we don't need extra caching
		// as we implemented our own caching solution, log unreturned err
		if err := registry.Destroy(); err != nil {
			e.Log.Errorf("failed to destroy registry: %w", err)
		}
	}()

	ref := image.SimpleReference(bundleImage)
	if err := registry.Pull(ctx, ref); err != nil {
		return err
	}

	if err := registry.Unpack(ctx, ref, tmpDirs["bundle"]); err != nil {
		return err
	}

	return e.ValidateBundle(registry, tmpDirs["bundle"])
}

func (e *DefaultBundleExtractor) ValidateBundle(registry *containerdregistry.Registry, tmpDir string) error {
	e.Log.Debugf("validating the unpacked bundle from %s", tmpDir)

	validator := opmbundle.NewImageValidator(registry, e.Log.(*logrus.Entry))
	if err := validator.ValidateBundleFormat(tmpDir); err != nil {
		return fmt.Errorf("bundle format validation failed: %w", err)
	}

	if err := validator.ValidateBundleContent(filepath.Join(tmpDir, opmbundle.ManifestsDir)); err != nil {
		return fmt.Errorf("bundle content validation failed: %w", err)
	}

	return nil
}

type tempDirs map[string]string

func createTempDirs() (tempDirs, error) {
	res := make(tempDirs)

	for _, dirName := range []string{"bundle", "containerd"} {
		tmpDir, err := ioutil.TempDir(os.TempDir(), dirName+"-")
		if err != nil {
			return nil, fmt.Errorf("unable to create %s tmpDir: %w", dirName, err)
		}
		res[dirName] = tmpDir
	}

	return res, nil
}

func (t tempDirs) CleanUp() error {
	var errs []string

	for _, dir := range t {
		if err := os.RemoveAll(dir); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ":"))
	}

	return nil
}
