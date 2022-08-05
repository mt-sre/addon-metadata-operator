package validator

import (
	"fmt"

	"github.com/operator-framework/operator-registry/pkg/registry"
	"golang.org/x/mod/semver"
)

func GetLatestBundle(bundles []*registry.Bundle) (*registry.Bundle, error) {
	if len(bundles) == 1 {
		return bundles[0], nil
	}

	latest := bundles[0]
	for _, bundle := range bundles[1:] {
		currVersion, err := getVersion(bundle)
		if err != nil {
			return nil, err
		}
		currLatestVersion, err := getVersion(latest)
		if err != nil {
			return nil, err
		}

		res := semver.Compare(currVersion, currLatestVersion)
		// If currVersion is greater than currLatestVersion
		if res == 1 {
			latest = bundle
		}
	}
	return latest, nil
}

func getVersion(bundle *registry.Bundle) (string, error) {
	csv, err := bundle.ClusterServiceVersion()
	if err != nil {
		return "", err
	}

	version, err := csv.GetVersion()
	if err != nil {
		return "", err
	}

	// Prefix a `v` infront of the version
	// so that semver package can parse it.
	return fmt.Sprintf("v%s", version), nil
}
