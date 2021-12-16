package validators

import (
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"golang.org/x/mod/semver"
)

func init() {
	Registry.Add(AM0003)
}

var AM0003 = types.Validator{
	Code:        "AM0003",
	Name:        "operator_name",
	Description: "Validate the operatorName matches csv.Name, csv.Replaces and bundle package annotation.",
	Runner:      validateOperatorName,
}

func validateOperatorName(mb types.MetaBundle) types.ValidatorResult {
	var failureMsgs []string
	operatorName := mb.AddonMeta.OperatorName

	for _, bundle := range mb.Bundles {
		b, err := newBundleDataAM0003(bundle)
		if err != nil {
			return Error(err)
		}

		if msg := checkCSVNameOrReplaces(b.csvName, operatorName); msg != "" {
			msg := fmt.Sprintf("bundle %s failed validation on csv.Name: %s.", b.nameVersion, msg)
			failureMsgs = append(failureMsgs, msg)
		}

		if b.csvReplaces != "" {
			if msg := checkCSVNameOrReplaces(b.csvReplaces, operatorName); msg != "" {
				msg := fmt.Sprintf("bundle '%s' failed validation on csv.Replaces: %s", b.nameVersion, msg)
				failureMsgs = append(failureMsgs, msg)
			}
		}

		if b.pkgNameAnnotation != operatorName {
			msg := fmt.Sprintf("bundle '%s' package annotation does not match operatorName '%s'", b.nameVersion, operatorName)
			failureMsgs = append(failureMsgs, msg)
		}

	}
	if len(failureMsgs) > 0 {
		return Fail(strings.Join(failureMsgs, ", "))
	}
	return Success()
}

func checkCSVNameOrReplaces(csvField, operatorName string) string {
	parts := strings.SplitN(csvField, ".", 2)
	if len(parts) != 2 {
		return fmt.Sprintf("could not split '%s' in two parts.", csvField)
	}
	bundleOperatorName := parts[0]
	bundleVersion := parts[1]

	if bundleOperatorName != operatorName {
		return fmt.Sprintf("invalid operatorName for '%s', should match '%s'", csvField, operatorName)
	}
	if !isValidSemver(bundleVersion) {
		return fmt.Sprintf("invalid semver '%s'", csvField)
	}
	return ""
}

func isValidSemver(version string) bool {
	if !strings.HasPrefix(version, "v") {
		version = fmt.Sprintf("v%s", version)
	}
	return semver.IsValid(version)
}

type bundleDataAM0003 struct {
	nameVersion       string
	csvName           string
	csvReplaces       string
	pkgNameAnnotation string
}

func newBundleDataAM0003(bundle *registry.Bundle) (bundleDataAM0003, error) {
	var bundleData bundleDataAM0003

	nameVersion, err := utils.GetBundleNameVersion(bundle)
	if err != nil {
		return bundleData, fmt.Errorf("could not get bundle name and version: %w", err)
	}

	csv, err := bundle.ClusterServiceVersion()
	if err != nil {
		return bundleData, fmt.Errorf("could not get csv for bundle '%s': %w", nameVersion, err)
	}

	replaces, err := csv.GetReplaces()
	if err != nil {
		return bundleData, fmt.Errorf("could not get csv.Replaces for bundle '%s': %w", nameVersion, err)
	}

	bundleData.csvName = csv.Name
	bundleData.csvReplaces = replaces
	bundleData.nameVersion = nameVersion
	bundleData.pkgNameAnnotation = bundle.Annotations.PackageName

	return bundleData, nil
}
