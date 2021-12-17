package validators

import (
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

// TODO - review failureMsg ++ missing test_validate_csv testing bundle
func ValidateCSVPresent(metabundle utils.MetaBundle) (bool, string, error) {
	if len(metabundle.Bundles) == 0 {
		return false, "No bundles present", nil
	}

	channels := metabundle.AddonMeta.Channels
	allOkay := true
	for _, channel := range *channels {
		requiredCsv := channel.CurrentCSV
		present, err := csvPresent(metabundle.Bundles, requiredCsv)
		if err != nil {
			return false, "", err
		}
		if !present {
			allOkay = false
		}
	}
	return allOkay, "", nil
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
