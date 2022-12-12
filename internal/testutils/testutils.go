package testutils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"

	"github.com/operator-framework/operator-registry/pkg/registry"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var rootDir string

func init() {
	_, b, _, _ := runtime.Caller(0)
	rootDir = filepath.Join(filepath.Dir(b), "..", "..")
}

/*
Implement a chaining directory interface. Examples:

	RootDir().TestData().MetadataV1().Legacy()
	RootDir().TestData().MetadataV2()
	RootDir().Bundles()
*/
type Tree string
type RootTree Tree
type TestDataTree Tree
type MetadataV1Tree Tree

func RootDir() RootTree {
	return RootTree(rootDir)
}

func (t RootTree) TestData() TestDataTree {
	return TestDataTree(filepath.Join(string(t), "internal", "testdata"))
}

func (t TestDataTree) MetadataV2() string {
	return filepath.Join(string(t), "metadata_v2")
}

func (t TestDataTree) Bundles() string {
	return filepath.Join(string(t), "bundles")
}

func (t TestDataTree) Validators() string {
	return filepath.Join(string(t), "validators")
}

func (t TestDataTree) MetadataV1() MetadataV1Tree {
	return MetadataV1Tree(filepath.Join(string(t), "metadata_v1"))
}

func (t MetadataV1Tree) Legacy() string {
	return filepath.Join(string(t), "legacy")
}

func (t MetadataV1Tree) ImageSets() string {
	return filepath.Join(string(t), "imagesets")
}

func RemoveDir(downloadDir string) {
	err := os.RemoveAll(downloadDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to cleanup download dir")
	}
}

func YamlToDynamicObj(yamlPath string) (unstructured.Unstructured, error) {
	pathTokens := strings.Split(yamlPath, ".")
	ext := pathTokens[len(pathTokens)-1]
	if ext != "yaml" && ext != "yml" {
		return unstructured.Unstructured{}, fmt.Errorf("non-yaml file found")
	}

	yamlAbsPath, err := filepath.Abs(yamlPath)
	if err != nil {
		return unstructured.Unstructured{}, err
	}
	yamlBytes, err := os.ReadFile(yamlAbsPath)
	if err != nil {
		return unstructured.Unstructured{}, err
	}
	var dynamicObj unstructured.Unstructured
	if err := yaml.Unmarshal(yamlBytes, &dynamicObj); err != nil {
		return unstructured.Unstructured{}, err
	}
	return dynamicObj, nil
}

func GetStringLiteralRef(s string) *string { return &s }

// DefaultSucceedingCandidates - returns a slice of valid metaBundles that are supposed
// to pass all validators successfully. If it is not the case, please make the required adjustments.
func DefaultSucceedingCandidates() ([]types.MetaBundle, error) {
	var res []types.MetaBundle
	refAddonStage, err := GetReferenceAddonStage()
	if err != nil {
		return nil, fmt.Errorf("Could not get reference-addon singleton, got %v.", err)
	}
	refAddonMetaBundle, err := refAddonStage.GetMetaBundle(*refAddonStage.MetaImageSet.ImageSetVersion)
	if err != nil {
		return nil, fmt.Errorf("Could not get reference-addon meta bundles, got %v.", err)
	}
	res = append(res, *refAddonMetaBundle)

	return res, nil
}

func NewBundle(name string, yamlFilePaths ...string) (*registry.Bundle, error) {
	var objs []*unstructured.Unstructured
	for _, path := range yamlFilePaths {
		obj, err := YamlToDynamicObj(path)
		if err != nil {
			return nil, fmt.Errorf("could not generate bundle: %v", err)
		}
		objs = append(objs, &obj)
	}
	return registry.NewBundle(name, &registry.Annotations{}, objs...), nil
}
