package validators

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func StringToPtr(s string) *string { return &s }

func yamlToDynamicObj(yamlPath string) (unstructured.Unstructured, error) {
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
	var csvObj unstructured.Unstructured
	if err := yaml.Unmarshal(yamlBytes, &csvObj); err != nil {
		return unstructured.Unstructured{}, err
	}
	return csvObj, nil
}
