package csvutils

import (
	"encoding/json"
	"strings"
	"unicode"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	operatorv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/registry"
	rbac "k8s.io/api/rbac/v1"
)

const wildCardStr = "*"

// Checks if secrets and configmaps without explicitly defined resource names
// are accessed at the cluster scope.
func CheckForConfidentialObjAccessAtClusterScope(csvPermissions *types.CSVPermissions) bool {
	filterConds := types.RuleFilter{
		PermissionType: types.ClusterPermissionType,
		Filters: []types.Filter{
			&types.APIGroupFilter{
				Params: types.FilterParams{
					Args:         []string{""},
					OperatorName: types.InOperator,
				},
			},
			&types.ResourcesFilter{
				Params: types.FilterParams{
					Args:         []string{"secrets", "configmaps"},
					OperatorName: types.AnyOperator,
				},
			},
			&types.ResourceNamesFilter{
				Params: types.FilterParams{
					Args:         []string{},
					OperatorName: types.DoesNotExistOperator,
				},
			},
		},
	}
	matchedRules := csvPermissions.FilterRules(filterConds)
	return len(matchedRules) > 0
}

// Checks if any rules have "*" defined in its apiGroup definition.
func WildCardApiGroupPresent(csvPermissions *types.CSVPermissions) bool {
	filterConds := types.RuleFilter{
		PermissionType: types.AllPermissionType,
		Filters: []types.Filter{
			&types.APIGroupFilter{
				Params: types.FilterParams{
					Args:         []string{wildCardStr},
					OperatorName: types.InOperator,
				},
			},
		},
	}
	matchedRules := csvPermissions.FilterRules(filterConds)
	return len(matchedRules) > 0
}

// Checks if any rules have "*" defined under resources.(For non-operator owned apis.)
func WildCardResourcePresent(csvPermissions *types.CSVPermissions, ownedApis []string) bool {
	filterConds := types.RuleFilter{
		PermissionType: types.AllPermissionType,
		Filters: []types.Filter{
			&types.APIGroupFilter{
				Params: types.FilterParams{
					Args:         ownedApis,
					OperatorName: types.NotEqualOperator,
				},
			},
			&types.ResourcesFilter{
				Params: types.FilterParams{
					Args:         []string{wildCardStr},
					OperatorName: types.InOperator,
				},
			},
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

func GetPermissions(csv *registry.ClusterServiceVersion) (*types.CSVPermissions, error) {
	var csvSpec operatorv1alpha1.ClusterServiceVersionSpec
	if err := json.Unmarshal(csv.Spec, &csvSpec); err != nil {
		return nil, err
	}
	clusterPermissions := csvSpec.InstallStrategy.StrategySpec.ClusterPermissions
	permissions := csvSpec.InstallStrategy.StrategySpec.Permissions

	return &types.CSVPermissions{
		ClusterPermissions: operatorPermissions2LocalPermissions(clusterPermissions),
		Permissions:        operatorPermissions2LocalPermissions(permissions),
	}, nil
}

func operatorPermissions2LocalPermissions(permissions []operatorv1alpha1.StrategyDeploymentPermissions) []types.Permission {
	res := make([]types.Permission, len(permissions))
	for _, permission := range permissions {
		res = append(res, types.Permission{
			ServiceAccountName: permission.ServiceAccountName,
			Rules:              operatorRule2localRules(permission.Rules),
		})
	}
	return res
}

func operatorRule2localRules(input []rbac.PolicyRule) []types.Rule {
	res := make([]types.Rule, len(input))
	for _, rule := range input {
		res = append(res, types.Rule{
			PolicyRule: rule,
		})
	}
	return res
}

func trimWhiteSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
