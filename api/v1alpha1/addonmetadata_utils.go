package v1alpha1

import (
	"encoding/json"

	ocmv1 "github.com/mt-sre/addon-metadata-operator/pkg/ocm/v1"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// FromYAML - instantiates an AddonMetadataSpec struct from yaml data
func (a *AddonMetadataSpec) FromYAML(data []byte) error {
	return yamlutil.Unmarshal(data, a)
}

// ToJSON - marshal AddonMetadata to JSON
func (a *AddonMetadata) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a *AddonMetadata) ToYAML() ([]byte, error) {
	return yaml.Marshal(a)
}

// CombineWithImageSet - Returns a new AddonMetadataSpec combined with the
// related imageSet fields. Using deep copy to avoid overriding the existing CR.
func (a *AddonMetadataSpec) CombineWithImageSet(imageSet *AddonImageSetSpec) (*AddonMetadataSpec, error) {
	combined := a.DeepCopy()

	imageSetVersion, err := imageSet.GetSemver()
	if err != nil {
		return nil, err
	}

	combined.IndexImage = &imageSet.IndexImage
	// TODO - do we need this? overrides latest?
	combined.ImageSetVersion = &imageSetVersion

	if imageSet.AddOnParameters != nil {
		params := make([]ocmv1.AddOnParameter, len(*imageSet.AddOnParameters))
		copy(params, *imageSet.AddOnParameters)
		combined.AddOnParameters = &params
	}

	if imageSet.AddOnRequirements != nil {
		requirements := make([]ocmv1.AddOnRequirement, len(*imageSet.AddOnRequirements))
		copy(requirements, *imageSet.AddOnRequirements)
		combined.AddOnRequirements = &requirements
	}

	if imageSet.SubOperators != nil {
		subOperators := make([]ocmv1.AddOnSubOperator, len(*imageSet.SubOperators))
		copy(subOperators, *imageSet.SubOperators)
		combined.SubOperators = &subOperators
	}

	return combined, nil
}
