package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterRules(t *testing.T) {
	inputRbacRules := CsvPermissions{
		ClusterPermissions: []Permissions{
			{
				Rules: []Rule{
					{
						name:            "rule-1",
						ApiGroups:       []string{"api_group_1"},
						Resources:       []string{"resource_1"},
						Verbs:           []string{"*"},
						ResourceNames:   []string{"sample1"},
						NonResourceURLs: []string{},
					},
					{
						name:            "rule-2",
						ApiGroups:       []string{"api_group_2", "api_group_3", "api_group_1"},
						Resources:       []string{"resource_2", "resource_3", "*"},
						Verbs:           []string{"*"},
						ResourceNames:   []string{"sample2", "sample3"},
						NonResourceURLs: []string{},
					},
				},
			},
		},
		Permissions: []Permissions{
			{
				Rules: []Rule{
					{
						name:            "rule-3",
						ApiGroups:       []string{"api_group_4"},
						Resources:       []string{"resource_5"},
						Verbs:           []string{"*"},
						ResourceNames:   []string{"sample4"},
						NonResourceURLs: []string{"port-forward"},
					},
				},
			},
		},
	}
	filterCases := []struct {
		input          RuleFilter
		expectedOutput []string
	}{
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				ApiGroupFilterObj: &FilterObj{
					Args:         []string{"api_group_1"},
					OperatorName: InOperator,
				},
			},
			expectedOutput: []string{"rule-1", "rule-2"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				ApiGroupFilterObj: &FilterObj{
					Args:         []string{"api_group_10"},
					OperatorName: InOperator,
				},
			},
			expectedOutput: []string{},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				ApiGroupFilterObj: &FilterObj{
					Args:         []string{"api_group_1"},
					OperatorName: NotInOperator,
				},
			},
			expectedOutput: []string{"rule-3"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				ResourceNamesFilterObj: &FilterObj{
					Args:         []string{"sample1"},
					OperatorName: InOperator,
				},
			},
			expectedOutput: []string{"rule-1"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				ApiGroupFilterObj: &FilterObj{
					Args:         []string{"api_group_1"},
					OperatorName: NotEqualOperator,
				},
			},
			expectedOutput: []string{"rule-2", "rule-3"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				ApiGroupFilterObj: &FilterObj{
					Args:         []string{"api_group_1", "api_group_3", "api_group_2"},
					OperatorName: EqualsOperator,
				},
			},
			expectedOutput: []string{"rule-2"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				ApiGroupFilterObj: &FilterObj{
					Args:         []string{"api_group_1", "api_group_3", "api_group_2"},
					OperatorName: EqualsOperator,
				},
			},
			expectedOutput: []string{"rule-2"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				NonResourceURLsFilterObj: &FilterObj{
					Args:         []string{},
					OperatorName: ExistsOperator,
				},
			},
			expectedOutput: []string{"rule-3"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				NonResourceURLsFilterObj: &FilterObj{
					Args:         []string{},
					OperatorName: DoesNotExistOperator,
				},
			},
			expectedOutput: []string{"rule-1", "rule-2"},
		},
		{
			input: RuleFilter{
				PermissionType: NameSpacedPermissionType,
				ApiGroupFilterObj: &FilterObj{
					Args:         []string{"api_group_1"},
					OperatorName: InOperator,
				},
			},
			expectedOutput: []string{},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				ApiGroupFilterObj: &FilterObj{
					Args:         []string{"api_group_1"},
					OperatorName: InOperator,
				},
				ResourcesFilterObj: &FilterObj{
					Args:         []string{"*"},
					OperatorName: InOperator,
				},
			},
			expectedOutput: []string{"rule-2"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				ApiGroupFilterObj: &FilterObj{
					Args:         []string{"api_group_1"},
					OperatorName: NotEqualOperator,
				},
				VerbsFilterObj: &FilterObj{
					Args:         []string{"*"},
					OperatorName: InOperator,
				},
			},
			expectedOutput: []string{"rule-2", "rule-3"},
		},
	}

	for _, testCase := range filterCases {
		res := inputRbacRules.FilterRules(testCase.input)
		resRuleNames := make([]string, 0)
		for _, item := range res {
			resRuleNames = append(resRuleNames, item.name)
		}
		require.Equal(t, testCase.expectedOutput, resRuleNames)
	}
}
