package validators

import (
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

// ValidateAddonLabel validates whether the 'label' field under an addon.yaml follows the format 'api.openshift.com/addon-<id>'
func ValidateCSVPresent(metabundle *utils.MetaBundle) (bool, error) {
	if len(metabundle.Bundles) == 0 {
		return false, fmt.Errorf("No bundles present!")
	}

	channels := metabundle.AddonMeta.Channels
	allOkay := true
	for _, channel := range channels {
		requiredCsv := channel.CurrentCSV
		present, err := csvPresent(metabundle.Bundles, requiredCsv)
		if err != nil {
			return false, err
		}
		if !present {
			allOkay = false
		}
	}
	return allOkay, nil
}

func csvPresent(bundles []registry.Bundle, requiredCsv string) (bool, error) {
	found := false
	for _, bundle := range bundles {
		csv, err := bundle.ClusterServiceVersion()
		if err != nil {
			return false, err
		}
		if csv.Name == requiredCsv {
			found = true
			return found, nil
		}
	}
	return found, nil
}
