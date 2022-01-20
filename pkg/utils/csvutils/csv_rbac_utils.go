package csvutils

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

const wildCardStr = "*"

// Checks if secrets and configmaps without explicitly defined resource names
// are accessed at the cluster scope.
func CheckForConfidentialObjAccessAtClusterScope(csvPermissions *types.CsvPermissions) bool {
	filterConds := types.RuleFilter{
		PermissionType: types.ClusterPermissionType,
		ApiGroupFilterObj: &types.FilterObj{
			Args:         []string{""},
			OperatorName: types.InOperator,
		},
		ResourcesFilterObj: &types.FilterObj{
			Args:         []string{"secrets", "configmaps"},
			OperatorName: types.AnyOperator,
		},
		ResourceNamesFilterObj: &types.FilterObj{
			Args:         []string{},
			OperatorName: types.DoesNotExistOperator,
		},
	}
	matchedRules := csvPermissions.FilterRules(filterConds)
	return len(matchedRules) > 0
}

// Checks if any rules have "*" defined in its apiGroup definition.
func WildCardApiGroupPresent(csvPermissions *types.CsvPermissions) bool {
	filterConds := types.RuleFilter{
		PermissionType: types.AllPermissionType,
		ApiGroupFilterObj: &types.FilterObj{
			Args:         []string{wildCardStr},
			OperatorName: types.InOperator,
		},
	}
	matchedRules := csvPermissions.FilterRules(filterConds)
	return len(matchedRules) > 0
}

// Checks if any rules have "*" defined under resources.(For non-operator owned apis.)
func WildCardResourcePresent(csvPermissions *types.CsvPermissions, ownedApis []string) bool {
	filterConds := types.RuleFilter{
		PermissionType: types.AllPermissionType,
		ApiGroupFilterObj: &types.FilterObj{
			Args:         ownedApis,
			OperatorName: types.NotEqualOperator,
		},
		ResourcesFilterObj: &types.FilterObj{
			Args:         []string{wildCardStr},
			OperatorName: types.InOperator,
		},
	}
	matchedRules := csvPermissions.FilterRules(filterConds)
	return len(matchedRules) > 0
}

func GetApisOwned(csv *registry.ClusterServiceVersion) ([]string, error) {
	ownedApis, _, err := csv.GetCustomResourceDefintions()
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, api := range ownedApis {
		if api != nil {
			result = append(result, trimWhiteSpace(api.Group))
		}
	}
	return result, nil
}

func GetPermissions(csv *registry.ClusterServiceVersion) (*types.CsvPermissions, error) {
	var objmap map[string]*json.RawMessage
	if err := json.Unmarshal(csv.Spec, &objmap); err != nil {
		return nil, err
	}
	installData, ok := objmap["install"]
	if !ok {
		return nil, fmt.Errorf("Failed to parse install spec from CSV")
	}
	var installMap map[string]*json.RawMessage
	if err := json.Unmarshal(*installData, &installMap); err != nil {
		return nil, err
	}

	installSpecBytes, ok := installMap["spec"]
	if !ok {
		return nil, fmt.Errorf("Failed to parse install spec from CSV")
	}
	csvPermissions := &types.CsvPermissions{}
	if err := json.Unmarshal(*installSpecBytes, csvPermissions); err != nil {
		return nil, fmt.Errorf("Failed to parse install spec from CSV")
	}

	return csvPermissions, nil
}

func trimWhiteSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
