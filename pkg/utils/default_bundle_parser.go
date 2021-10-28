package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/operator-framework/operator-registry/pkg/registry"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

type DefaultBundleParser struct{}

func (obj DefaultBundleParser) ParseBundles(addonName, manifestsDir string) ([]registry.Bundle, error) {
	var bundles []registry.Bundle
	bundlesDir, err := os.ReadDir(manifestsDir)
	if err != nil {
		return []registry.Bundle{}, err
	}
	for _, bundlePath := range bundlesDir {
		if bundlePath.IsDir() {
			bundle := registry.Bundle{
				Name:        addonName,
				Annotations: &registry.Annotations{},
			}
			k8sObjs := parseK8sObjects(bundlePath, addonName, manifestsDir)
			for _, k8sobj := range k8sObjs {
				bundle.Add(&k8sobj)
			}
			bundles = append(bundles, bundle)
		}
	}
	return bundles, nil
}

func parseK8sObjects(bundleDir os.DirEntry, addonName, manifestsDir string) []unstructured.Unstructured {
	var result []unstructured.Unstructured
	dirPath := filepath.Join(manifestsDir, bundleDir.Name())
	items, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Warnf("Failed to read %s bundle directory!\n", bundleDir.Name())
	}

	for _, item := range items {
		filePath := filepath.Join(dirPath, item.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Warnf("Failed to read manifest files in %s bundle!\n", bundleDir.Name())
		}
		dec := k8syaml.NewYAMLOrJSONDecoder(strings.NewReader(string(data)), 30)
		k8sFile := unstructured.Unstructured{}
		err = dec.Decode(&k8sFile)
		if err != nil {
			log.Warnf("Failed to parse k8s manifest file in %s bundle\n", bundleDir.Name())
		} else {
			result = append(result, k8sFile)
		}
	}
	return result
}
