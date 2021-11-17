package v1alpha1

import "k8s.io/apimachinery/pkg/util/yaml"

// FromYAML - instantiates an AddonMetadataSpec struct from yaml data
func (a *AddonMetadataSpec) FromYAML(data []byte) error {
	return yaml.Unmarshal(data, a)
}

// FromYAML - instantiates an AddonImageSetSpec struct from yaml data
func (a *AddonImageSetSpec) FromYAML(data []byte) error {
	return yaml.Unmarshal(data, a)
}
