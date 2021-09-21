package v1alpha1

import "k8s.io/apimachinery/pkg/util/yaml"

// *****
// Helper types
// *****
// Channel - list all channels for a given operator
type Channel struct {
	Name       string `json:"name"`
	CurrentCSV string `json:"currentCSV"`
}

// *****
// Helper funcs/methods
// *****
// FromYAML - instantiates an AddonMetadata struct from yaml data
func (a *AddonMetadata) FromYAML(data []byte) error {
	return yaml.Unmarshal(data, a)
}
