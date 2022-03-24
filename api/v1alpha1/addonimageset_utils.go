package v1alpha1

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// FromYAML - instantiates an AddonImageSetSpec struct from yaml data
func (a *AddonImageSetSpec) FromYAML(data []byte) error {
	return yaml.Unmarshal(data, a)
}

// ToJSON - marshal ojbect as JSON
func (a *AddonImageSet) ToJSON() ([]byte, error) {
	return json.Marshal(a)
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
