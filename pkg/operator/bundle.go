package operator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/blang/semver/v4"
	opsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	opmbundle "github.com/operator-framework/operator-registry/pkg/lib/bundle"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

func NewBundleFromDirectory(path string) (Bundle, error) {
	unstObjs, err := readAllManifests(filepath.Join(path, opmbundle.ManifestsDir))
	if err != nil {
		return Bundle{}, fmt.Errorf("reading manifests: %w", err)
	}

	annotations, err := readAnnotations(filepath.Join(path, opmbundle.MetadataDir))
	if err != nil {
		return Bundle{}, fmt.Errorf("reading annotations: %w", err)
	}

	regBundle := registry.NewBundle(annotations.PackageName, annotations, unstObjs...)
	bundle, err := NewBundleFromRegistryBundle(*regBundle)
	if err != nil {
		return Bundle{}, fmt.Errorf("generating bundle: %w", err)
	}

	return bundle, nil
}

func readAllManifests(manifestsDir string) ([]*unstructured.Unstructured, error) {
	var objs []*unstructured.Unstructured

	items, err := os.ReadDir(manifestsDir)
	if err != nil {
		return nil, fmt.Errorf("reading manifests dir: %w", err)
	}

	for _, item := range items {
		path := filepath.Join(manifestsDir, item.Name())

		manifest, err := readManifest(path)
		if err != nil {
			return nil, fmt.Errorf("reading manifest %q: %w", path, err)
		}

		objs = append(objs, manifest)
	}
	return objs, nil
}

func readManifest(path string) (*unstructured.Unstructured, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	dec := k8syaml.NewYAMLOrJSONDecoder(strings.NewReader(string(data)), 30)

	var manifest unstructured.Unstructured

	if err = dec.Decode(&manifest); err != nil {
		return nil, fmt.Errorf("decoding manifest: %w", err)
	}

	return &manifest, nil
}

func readAnnotations(metadataDir string) (*registry.Annotations, error) {
	path := filepath.Join(metadataDir, opmbundle.AnnotationsFile)

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", path, err)
	}

	var annotationsFile registry.AnnotationsFile
	if err := yaml.Unmarshal(content, &annotationsFile); err != nil {
		return nil, fmt.Errorf("unmarshalling file %q: %w", path, err)
	}

	return &annotationsFile.Annotations, nil
}

func NewBundleFromRegistryBundle(b registry.Bundle) (Bundle, error) {
	var annotations Annotations
	if as := b.Annotations; as != nil {
		annotations = NewAnnotationsFromRegistryAnnotations(*as)
	}

	ver, err := b.Version()
	if err != nil {
		return Bundle{}, fmt.Errorf("getting bundle version: %w", err)
	}

	regCSV, err := b.ClusterServiceVersion()
	if err != nil {
		return Bundle{}, fmt.Errorf("getting bundle ClusterServiceVersion: %w", err)
	}

	var csv ClusterServiceVersion

	if regCSV != nil {
		csv, err = NewClusterServiceVersionfromRegistryCSV(*regCSV)
		if err != nil {
			return Bundle{}, fmt.Errorf("setting ClusterServiceVersion values: %w", err)
		}
	}

	return Bundle{
		Annotations:           annotations,
		Name:                  b.Name,
		Package:               b.Package,
		Channels:              b.Channels,
		BundleImage:           b.BundleImage,
		Version:               ver,
		ClusterServiceVersion: csv,
	}, nil
}

type Bundle struct {
	Annotations           Annotations
	ClusterServiceVersion ClusterServiceVersion
	BundleImage           string
	Channels              []string
	Name                  string
	Package               string
	Version               string
}

func (b *Bundle) GetNameVersion() string {
	return fmt.Sprintf("%s:%s", b.Name, b.Version)
}

func NewAnnotationsFromRegistryAnnotations(as registry.Annotations) Annotations {
	return Annotations{
		PackageName:        as.PackageName,
		Channels:           strings.Split(as.Channels, ","),
		DefaultChannelName: as.DefaultChannelName,
	}
}

type Annotations struct {
	PackageName        string
	Channels           []string
	DefaultChannelName string
}

func NewClusterServiceVersionfromRegistryCSV(csv registry.ClusterServiceVersion) (ClusterServiceVersion, error) {
	var spec opsv1alpha1.ClusterServiceVersionSpec
	if err := json.Unmarshal(csv.Spec, &spec); err != nil {
		return ClusterServiceVersion{}, fmt.Errorf("unmarshalling csv spec: %w", err)
	}

	owned, required, err := csv.GetCustomResourceDefintions()
	if err != nil {
		return ClusterServiceVersion{}, fmt.Errorf("getting CustomResourceDefinitions: %w", err)
	}

	ownedCRDs := make([]CustomResourceDefinition, 0, len(owned))
	for _, key := range owned {
		if key == nil {
			continue
		}

		ownedCRDs = append(ownedCRDs, NewCustomeResourceDefinitionFromRegistryDefinitionKey(*key))
	}

	requiredCRDs := make([]CustomResourceDefinition, 0, len(required))
	for _, key := range required {
		if key == nil {
			continue
		}

		requiredCRDs = append(requiredCRDs, NewCustomeResourceDefinitionFromRegistryDefinitionKey(*key))
	}

	return ClusterServiceVersion{
		Name:                              csv.Name,
		OwnedCustomResourceDefinitions:    ownedCRDs,
		RequiredCustomResourceDefinitions: requiredCRDs,
		Spec:                              spec,
	}, nil
}

type ClusterServiceVersion struct {
	Name                              string
	OwnedCustomResourceDefinitions    []CustomResourceDefinition
	RequiredCustomResourceDefinitions []CustomResourceDefinition
	Spec                              opsv1alpha1.ClusterServiceVersionSpec
}

func NewCustomeResourceDefinitionFromRegistryDefinitionKey(key registry.DefinitionKey) CustomResourceDefinition {
	return CustomResourceDefinition{
		Name:    key.Name,
		Group:   key.Group,
		Kind:    key.Kind,
		Version: key.Version,
	}
}

type CustomResourceDefinition struct {
	Name                 string
	Group, Kind, Version string
}

func HeadBundle(bundles ...Bundle) (Bundle, bool) {
	if len(bundles) < 1 {
		return Bundle{}, false
	}

	if len(bundles) == 1 {
		return bundles[0], true
	}

	ordered := OrderedBundles(bundles)

	sort.Sort(ordered)

	return ordered[0], true
}

type OrderedBundles []Bundle

func (l OrderedBundles) Len() int      { return len(l) }
func (l OrderedBundles) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l OrderedBundles) Less(i, j int) bool {
	iVer, _ := semver.ParseTolerant(l[i].Version)
	jVer, _ := semver.ParseTolerant(l[j].Version)

	return jVer.LE(iVer) // reverse sort
}
