package v1alpha1

import "k8s.io/apimachinery/pkg/util/yaml"

// FromYAML - instantiates an AddonMetadata struct from yaml data
func (a *AddonMetadata) FromYAML(data []byte) error {
	return yaml.Unmarshal(data, a)
}

// FromYAML - instantiates an AddonMetadata struct from yaml data
func (a *AddonMetadataSpec) FromYAML(data []byte) error {
	return yaml.Unmarshal(data, a)
}
