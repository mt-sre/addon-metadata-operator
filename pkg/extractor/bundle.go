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

	"github.com/operator-framework/operator-registry/pkg/image"
	"github.com/operator-framework/operator-registry/pkg/image/containerdregistry"
	opmbundle "github.com/operator-framework/operator-registry/pkg/lib/bundle"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

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
		extractor.Cache = NewBundleMemoryCache(NewJSONSnappyEncoder())
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

func (e *DefaultBundleExtractor) Extract(ctx context.Context, bundleImage string) (*registry.Bundle, error) {
	bundle, err := e.Cache.Get(bundleImage)
	if err != nil {
		return nil, fmt.Errorf("cache error: %w", err)
	}
	if bundle != nil {
		e.Log.Debugf("cache hit for '%s'", bundleImage)
		return bundle, nil
	}

	e.Log.Debugf("cache miss for '%s'", bundleImage)
	tmpDirs, err := createTempDirs()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tmpDirs.CleanUp(); err != nil {
			e.Log.Errorf("failed to cleanup tmpDirs: %w", err)
		}
	}()

	if err := e.unpackAndValidateBundle(ctx, bundleImage, tmpDirs); err != nil {
		return nil, fmt.Errorf("failed to unpack and validate bundle: %w", err)
	}

	bundle, err = e.loadBundle(tmpDirs["bundle"])
	if err != nil {
		return nil, err
	}
	bundle.BundleImage = bundleImage // not set by OPM

	if err := e.Cache.Set(bundleImage, bundle); err != nil {
		return nil, fmt.Errorf("unable to cache bundle: %w", err)
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
		containerdregistry.SkipTLS(false),
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

	return e.validateBundle(registry, tmpDirs["bundle"])
}

func (e *DefaultBundleExtractor) validateBundle(registry *containerdregistry.Registry, tmpDir string) error {
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

func (e *DefaultBundleExtractor) loadBundle(tmpDir string) (*registry.Bundle, error) {
	e.Log.Debugf("loading the bundle from tmpDir: %s", tmpDir)
	unstObjs, err := readAllManifests(filepath.Join(tmpDir, opmbundle.ManifestsDir))
	if err != nil {
		return nil, err
	}
	annotations, err := readAnnotations(filepath.Join(tmpDir, opmbundle.MetadataDir))
	if err != nil {
		return nil, err
	}
	return registry.NewBundle(annotations.PackageName, annotations, unstObjs...), nil
}

func readAllManifests(manifestsDir string) ([]*unstructured.Unstructured, error) {
	unstObjs := []*unstructured.Unstructured{}
	items, err := ioutil.ReadDir(manifestsDir)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		path := filepath.Join(manifestsDir, item.Name())
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("unable to read file %s, got %w", path, err)
		}

		dec := k8syaml.NewYAMLOrJSONDecoder(strings.NewReader(string(data)), 30)
		k8sFile := &unstructured.Unstructured{}
		if err = dec.Decode(k8sFile); err != nil {
			return nil, fmt.Errorf("unable to decode file %s, got %w", path, err)
		}

		unstObjs = append(unstObjs, k8sFile)
	}
	return unstObjs, nil
}

func readAnnotations(metadataDir string) (*registry.Annotations, error) {
	path := filepath.Join(metadataDir, opmbundle.AnnotationsFile)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read file '%s': %w", path, err)
	}

	var annotationsFile registry.AnnotationsFile
	if err = yaml.Unmarshal(content, &annotationsFile); err != nil {
		return nil, fmt.Errorf("unable to unmarshal file '%s': %w", path, err)
	}

	return &annotationsFile.Annotations, nil
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
