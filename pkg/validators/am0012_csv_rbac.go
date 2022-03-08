package validators

import (
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils/csvutils"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"golang.org/x/mod/semver"
)

func init() {
	Registry.Add(AM0012)
}

var AM0012 = types.Validator{
	Code:        "AM0012",
	Name:        "csv_permissions",
	Description: "validates the permissions specified in the csv",
	Runner:      validateCsvRbac,
}

func validateCsvRbac(mb types.MetaBundle) types.ValidatorResult {
	// Only run validations on the latest bundle
	latestBundle, err := getLatestBundle(mb.Bundles)
	if err != nil {
		return Error(err)
	}

	csv, err := latestBundle.ClusterServiceVersion()
	if err != nil {
		return Error(err)
	}
	apisOwnedByOperator, err := csvutils.GetApisOwned(csv)
	if err != nil {
		return Error(err)
	}

	permissions, err := csvutils.GetPermissions(csv)
	if err != nil {
		return Error(err)
	}

	validationErrors := []string{}
	validationErrors, err = validateApiGroups(permissions, validationErrors)
	if err != nil {
		return Error(err)
	}
	validationErrors, err = validateWildcardInResources(permissions, apisOwnedByOperator, validationErrors)
	if err != nil {
		return Error(err)
	}
	validationErrors, err = validateConfidentialObjAccessAtClusterScope(permissions, validationErrors)
	if err != nil {
		return Error(err)
	}

	if len(validationErrors) == 0 {
		return Success()
	}
	failureMsg := "CSV rbac validation errors: \n" + strings.Join(validationErrors, "\n")
	return Fail(failureMsg)
}

func validateApiGroups(permissions *types.CSVPermissions, existingValidationErrs []string) ([]string, error) {
	if csvutils.WildCardApiGroupPresent(permissions) {
		errorMsg := "Wild card string used under api group/s"
		return append(existingValidationErrs, errorMsg), nil
	}
	return existingValidationErrs, nil
}

func validateWildcardInResources(csvPermissions *types.CSVPermissions,
	apisOwnedByOperator []string,
	existingValidationErrs []string) ([]string, error) {
	if csvutils.WildCardResourcePresent(csvPermissions, apisOwnedByOperator) {
		errorMsg := "Wild card string used under resource/s not owned by the operator"
		return append(existingValidationErrs, errorMsg), nil
	}
	return existingValidationErrs, nil
}

func validateConfidentialObjAccessAtClusterScope(csvPermissions *types.CSVPermissions,
	existingValidationErrs []string) ([]string, error) {
	if csvutils.CheckForConfidentialObjAccessAtClusterScope(csvPermissions) {
		errorMsg := "config maps/Secrets access rules present at the cluster scope"
		return append(existingValidationErrs, errorMsg), nil
	}
	return existingValidationErrs, nil
}

func getLatestBundle(bundles []registry.Bundle) (*registry.Bundle, error) {
	if len(bundles) == 1 {
		return &bundles[0], nil
	}

	latest := bundles[0]
	for _, bundle := range bundles[1:] {
		currVersion, err := getVersion(&bundle)
		if err != nil {
			return nil, err
		}
		currLatestVersion, err := getVersion(&latest)
		if err != nil {
			return nil, err
		}

		res := semver.Compare(currVersion, currLatestVersion)
		// If currVersion is greater than currLatestVersion
		if res == 1 {
			latest = bundle
		}
	}
	return &latest, nil
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
