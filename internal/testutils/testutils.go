package testutils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"io/ioutil"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"

	"github.com/operator-framework/operator-registry/pkg/registry"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var rootDir string

func init() {
	_, b, _, _ := runtime.Caller(0)
	rootDir = path.Join(filepath.Dir(b), "..", "..")
}

func AddonsIndexImageDir() string {
	return path.Join(TestdataDir(), "addons-indeximage")
}

func AddonsImagesetDir() string {
	return path.Join(TestdataDir(), "addons-imageset")
}

func TestdataDir() string {
	return path.Join(RootDir(), "internal", "testdata")
}

func RootDir() string {
	return rootDir
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
	yamlBytes, err := ioutil.ReadFile(yamlAbsPath)
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

func NewBundle(name string, yamlFilePaths ...string) (registry.Bundle, error) {
	var objs []*unstructured.Unstructured
	for _, path := range yamlFilePaths {
		obj, err := YamlToDynamicObj(path)
		if err != nil {
			return registry.Bundle{}, fmt.Errorf("could not generate bundle: %v", err)
		}
		objs = append(objs, &obj)
	}
	return *registry.NewBundle(name, &registry.Annotations{}, objs...), nil
}

func EmptyMetaBundle() types.MetaBundle {
	return types.MetaBundle{
		AddonMeta: &v1alpha1.AddonMetadataSpec{
			ID: "random-operator",
		},
	}
}
