package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/operator-framework/operator-registry/pkg/registry"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

type DefaultBundleParser struct{}

func (obj DefaultBundleParser) ParseBundles(addonName, manifestsDir string) ([]registry.Bundle, error) {
	var bundles []registry.Bundle
	var allErrs []error

	if !checkFileExists(manifestsDir) {
		return []registry.Bundle{}, fmt.Errorf(
			"can't find any bundles for the operator '%s'. Looked in: '%s'",
			addonName,
			manifestsDir,
		)
	}
	bundlesDir, err := os.ReadDir(manifestsDir)
	if err != nil {
		return []registry.Bundle{}, err
	}
	for _, bundlePath := range bundlesDir {
		if bundlePath.IsDir() {
			// TODO - (sblaisdo) enable after we migrate extraction format to bundles instead of packageManifest
			// annotations, err := readAnnotations(filepath.Join(manifestsDir, bundlePath.Name()))
			// if err != nil {
			//	allErrs = append(allErrs, err)
			// }
			bundle := registry.Bundle{
				Name:        addonName,
				Annotations: &registry.Annotations{},
			}
			k8sObjs, err := parseK8sObjects(bundlePath, addonName, manifestsDir)
			if err != nil {
				allErrs = append(allErrs, err)
			} else {
				for index := range k8sObjs {
					bundle.Add(&k8sObjs[index])
				}
				bundles = append(bundles, bundle)
			}
		}
	}
	if len(allErrs) != 0 {
		errMsgPrefix := "allErrs while parsing bundles: \n"
		return bundles, concatParseErrors(allErrs, errMsgPrefix)
	}

	if len(bundles) == 0 {
		return []registry.Bundle{}, errors.New("No bundles were found.")
	}
	return bundles, nil
}

func parseK8sObjects(bundleDir os.DirEntry, addonName, manifestsDir string) ([]unstructured.Unstructured, error) {
	var result []unstructured.Unstructured
	dirPath := filepath.Join(manifestsDir, bundleDir.Name())
	items, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return []unstructured.Unstructured{}, err
	}
	var errs []error
	for _, item := range items {
		filePath := filepath.Join(dirPath, item.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			err = fmt.Errorf("Failed to read manifest files in %s bundle!. Error: %s", bundleDir.Name(), err.Error())
			errs = append(errs, err)
			continue
		}
		dec := k8syaml.NewYAMLOrJSONDecoder(strings.NewReader(string(data)), 30)
		k8sFile := unstructured.Unstructured{}
		err = dec.Decode(&k8sFile)
		if err != nil {
			err = fmt.Errorf("Failed to parse k8s manifest file in %s bundle. Error: %s", bundleDir.Name(), err.Error())
			errs = append(errs, err)
		} else {
			result = append(result, k8sFile)
		}
	}
	if len(errs) != 0 {
		bundleIdentifier := strings.Join([]string{addonName, bundleDir.Name()}, "/")
		errorMsgPrefix := fmt.Sprintf(
			"Parsing k8s objects for the bundle: %s failed with the following errors:",
			bundleIdentifier,
		)
		return result, concatParseErrors(errs, errorMsgPrefix)
	}
	return result, nil
}

func concatParseErrors(errs []error, errMsgPrefix string) error {
	if len(errs) == 0 {
		return nil
	}
	errStrs := make([]string, len(errs))
	for _, err := range errs {
		errStrs = append(errStrs, err.Error())
	}
	errStrsJoined := strings.Join(errStrs, "\n")
	return fmt.Errorf(
		"\n %s",
		errStrsJoined,
	)
}

// TODO - (sblaisdo) enable when extration format is bundles not packageManifest
// func readAnnotations(bundlePath string) (*registry.Annotations, error) {
//	annotationsPath := filepath.Join(bundlePath, "metadata", "annotations.yaml")
//	content, err := ioutil.ReadFile(annotationsPath)
//	if err != nil {
//		return nil, fmt.Errorf("Could not read annotationsPath %v, got %v.", annotationsPath, err)
//	}

//	var annotationsFile registry.AnnotationsFile
//	err = yaml.Unmarshal(content, &annotationsFile)
//	if err != nil {
//		return nil, fmt.Errorf("Could not unmarshal annotationsFile, got %v.", err)
//	}

//	return &annotationsFile.Annotations, nil
// }
