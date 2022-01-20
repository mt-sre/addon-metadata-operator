package types

import (
	"fmt"
	"sort"
)

type CsvPermissions struct {
	ClusterPermissions []Permissions `json:"clusterPermissions"`
	Permissions        []Permissions `json:"permissions"`
}

type Permissions struct {
	ServiceAccount string `json:"serviceAccount"`
	Rules          []Rule `json:"rules"`
}

type Rule struct {
	name            string   // used in tests
	ApiGroups       []string `json:"apiGroups"`
	Resources       []string `json:"resources"`
	Verbs           []string `json:"verbs"`
	ResourceNames   []string `json:"resourceNames"`
	NonResourceURLs []string `json:"nonResourceURLs"`
}

type operator string

type permissionType string

type RuleFilter struct {
	PermissionType           permissionType
	ApiGroupFilterObj        *FilterObj
	ResourcesFilterObj       *FilterObj
	VerbsFilterObj           *FilterObj
	ResourceNamesFilterObj   *FilterObj
	NonResourceURLsFilterObj *FilterObj
}

type FilterObj struct {
	Args         []string
	OperatorName operator
}

type filterFunc func(*Rule, RuleFilter) *Rule

const (
	apiGroupFilterType        = "apiGroupFilter"
	resourcesFilterType       = "resourcesFilter"
	verbsFilterType           = "verbsFilter"
	resourceNamesFilterType   = "resourceNamesFilter"
	nonResourceURLsFilterType = "nonResourceURLsFilter"
)

var (
	InOperator           operator = "IN"
	NotInOperator        operator = "NOT_IN"
	EqualsOperator       operator = "EQUAL"
	NotEqualOperator     operator = "NOT_EQUAL"
	ExistsOperator       operator = "EXISTS"
	DoesNotExistOperator operator = "DOES_NOT_EXIST"
	AnyOperator          operator = "ANY"
)

var (
	AllPermissionType        permissionType = "all"
	NameSpacedPermissionType permissionType = "namespaced"
	ClusterPermissionType    permissionType = "clusterScoped"
)

func genFilter(attrFetcher func(*Rule) []string, filterType string) filterFunc {
	return func(rule *Rule, filter RuleFilter) *Rule {
		args := attrFetcher(rule)
		filterObj := func() *FilterObj {
			switch filterType {
			case apiGroupFilterType:
				return filter.ApiGroupFilterObj
			case resourcesFilterType:
				return filter.ResourcesFilterObj
			case verbsFilterType:
				return filter.VerbsFilterObj
			case resourceNamesFilterType:
				return filter.ResourceNamesFilterObj
			case nonResourceURLsFilterType:
				return filter.NonResourceURLsFilterObj
			default:
				panic("Unsupported filter type")
			}
		}()
		if eval(args, filterObj) {
			return rule
		}
		// Return nil if the rule doesnt match the filtering condition.
		return nil
	}
}

var (
	apiGroupFilter = genFilter(
		func(r *Rule) []string {
			if r != nil {
				return r.ApiGroups
			}
			return []string{}
		},
		apiGroupFilterType,
	)
	resourcesFilter = genFilter(
		func(r *Rule) []string {
			if r != nil {
				return r.Resources
			}
			return []string{}
		},
		resourcesFilterType,
	)
	verbsFilter = genFilter(
		func(r *Rule) []string {
			if r != nil {
				return r.Verbs
			}
			return []string{}
		},
		verbsFilterType,
	)
	resourceNamesFilter = genFilter(
		func(r *Rule) []string {
			if r != nil {
				return r.ResourceNames
			}
			return []string{}
		},
		resourceNamesFilterType,
	)
	nonResourceURLsFilter = genFilter(
		func(r *Rule) []string {
			if r != nil {
				return r.NonResourceURLs
			}
			return []string{}
		},
		nonResourceURLsFilterType,
	)
)

// Returns the list of rules matching the filtering conditions
func (cp CsvPermissions) FilterRules(ruleMatcher RuleFilter) []Rule {
	concernedPermissionRules := func() []Permissions {
		switch ruleMatcher.PermissionType {
		case AllPermissionType:
			return append(cp.ClusterPermissions, cp.Permissions...)
		case NameSpacedPermissionType:
			return cp.Permissions
		case ClusterPermissionType:
			return cp.ClusterPermissions
		default:
			return []Permissions{}
		}
	}()
	filteredRules := make([]Rule, 0)
	for _, permissionRule := range concernedPermissionRules {
		for _, rule := range permissionRule.Rules {
			res := runFilters(getAllAttributeFilters(), &rule, ruleMatcher)
			if res != nil {
				filteredRules = append(filteredRules, rule)
			}
		}
	}

	return filteredRules
}

func runFilters(filters []filterFunc, rule *Rule, condition RuleFilter) *Rule {
	if len(filters) == 0 || rule == nil {
		return rule
	}

	for _, filter := range filters {
		res := filter(rule, condition)
		if res == nil {
			return nil
		}
	}
	return rule
}

func getAllAttributeFilters() []filterFunc {
	return []filterFunc{
		apiGroupFilter,
		resourcesFilter,
		verbsFilter,
		resourceNamesFilter,
		nonResourceURLsFilter,
	}
}

func includes(items []string, itemsToBePresent []string) bool {
	itemsMap := sliceToSet(items)
	for _, item := range itemsToBePresent {
		if _, ok := itemsMap[item]; !ok {
			return false
		}
	}
	return true
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)

	for index := range a {
		if a[index] != b[index] {
			return false
		}
	}
	return true
}

// Checks if any element in b is present in a.
func any(a, b []string) bool {
	source := sliceToSet(a)
	for _, item := range b {
		if _, ok := source[item]; ok {
			return true
		}
	}
	return false
}

func sliceToSet(items []string) map[string]struct{} {
	res := make(map[string]struct{}, len(items))
	for _, item := range items {
		res[item] = struct{}{}
	}
	return res
}

func eval(ruleArgs []string, filterObj *FilterObj) bool {
	if filterObj == nil {
		return true
	}
	switch filterObj.OperatorName {
	case InOperator:
		return includes(ruleArgs, filterObj.Args)
	case NotInOperator:
		return !includes(ruleArgs, filterObj.Args)
	case EqualsOperator:
		return equal(ruleArgs, filterObj.Args)
	case NotEqualOperator:
		return !equal(ruleArgs, filterObj.Args)
	case ExistsOperator:
		return len(ruleArgs) > 0
	case DoesNotExistOperator:
		return len(ruleArgs) == 0
	case AnyOperator:
		return any(ruleArgs, filterObj.Args)
	default:
		panic(fmt.Sprintf("eval: Unsupported operator %s", filterObj.OperatorName))
	}
}
