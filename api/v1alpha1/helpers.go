package v1alpha1

import (
	"fmt"
	"strings"

	ocmv1 "github.com/mt-sre/addon-metadata-operator/pkg/ocm/v1"
	"golang.org/x/mod/semver"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// FromYAML - instantiates an AddonMetadataSpec struct from yaml data
func (a *AddonMetadataSpec) FromYAML(data []byte) error {
	return yaml.Unmarshal(data, a)
}

// FromYAML - instantiates an AddonImageSetSpec struct from yaml data
func (a *AddonImageSetSpec) FromYAML(data []byte) error {
	return yaml.Unmarshal(data, a)
}

// GetSemver - Returns the semver version matching "MAJOR.MINOR.PATCH".
func (a *AddonImageSetSpec) GetSemver() (string, error) {
	parts := strings.SplitN(a.Name, ".", 2)
	if len(parts) == 2 {
		version := parts[1]
		if semver.IsValid(version) {
			return strings.TrimPrefix(version, "v"), nil
		}
	}
	return "", fmt.Errorf("Could not parse the imageSet name as a valid semver, %v.", a.Name)
}

func (a *AddonMetadataSpec) PatchWithImageSet(imageSet *AddonImageSetSpec) error {
	imageSetVersion, err := imageSet.GetSemver()
	if err != nil {
		return err
	}
	a.IndexImage = &imageSet.IndexImage
	a.ImageSetVersion = &imageSetVersion

	// overwrite all metadata shared fields
	a.AddOnParameters = nil
	a.AddOnRequirements = nil
	a.SubOperators = nil
	if imageSet.AddOnParameters != nil {
		params := make([]ocmv1.AddOnParameter, len(*imageSet.AddOnParameters))
		copy(params, *imageSet.AddOnParameters)
		a.AddOnParameters = &params
	}
	if imageSet.AddOnRequirements != nil {
		requirements := make([]ocmv1.AddOnRequirement, len(*imageSet.AddOnRequirements))
		copy(requirements, *imageSet.AddOnRequirements)
		a.AddOnRequirements = &requirements
	}
	if imageSet.SubOperators != nil {
		subOperators := make([]ocmv1.AddOnSubOperator, len(*imageSet.SubOperators))
		copy(subOperators, *imageSet.SubOperators)
		a.SubOperators = &subOperators
	}
	return nil
}
