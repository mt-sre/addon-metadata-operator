package types

import (
	"testing"

	rbac "k8s.io/api/rbac/v1"

	"github.com/stretchr/testify/require"
)

func TestFilterRules(t *testing.T) {
	inputRbacRules := CSVPermissions{
		ClusterPermissions: []Permission{
			{
				Rules: []Rule{
					{
						name: "rule-1",
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{"api_group_1"},
							Resources:       []string{"resource_1"},
							Verbs:           []string{"*"},
							ResourceNames:   []string{"sample1"},
							NonResourceURLs: []string{},
						},
					},
					{
						name: "rule-2",
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{"api_group_2", "api_group_3", "api_group_1"},
							Resources:       []string{"resource_2", "resource_3", "*"},
							Verbs:           []string{"*"},
							ResourceNames:   []string{"sample2", "sample3"},
							NonResourceURLs: []string{},
						},
					},
				},
			},
		},
		Permissions: []Permission{
			{
				Rules: []Rule{
					{
						name: "rule-3",
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{"api_group_4"},
							Resources:       []string{"resource_5"},
							Verbs:           []string{"*"},
							ResourceNames:   []string{"sample4"},
							NonResourceURLs: []string{"port-forward"},
						},
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
				Filters: []Filter{
					&APIGroupFilter{
						Params: FilterParams{
							Args:         []string{"api_group_1"},
							OperatorName: InOperator,
						},
					},
				},
			},
			expectedOutput: []string{"rule-1", "rule-2"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				Filters: []Filter{
					&APIGroupFilter{
						Params: FilterParams{
							Args:         []string{"api_group_10"},
							OperatorName: InOperator,
						},
					},
				},
			},
			expectedOutput: []string{},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				Filters: []Filter{
					&APIGroupFilter{
						Params: FilterParams{
							Args:         []string{"api_group_1"},
							OperatorName: NotInOperator,
						},
					},
				},
			},
			expectedOutput: []string{"rule-3"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				Filters: []Filter{
					&ResourceNamesFilter{
						Params: FilterParams{
							Args:         []string{"sample1"},
							OperatorName: InOperator,
						},
					},
				},
			},
			expectedOutput: []string{"rule-1"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				Filters: []Filter{
					&APIGroupFilter{
						Params: FilterParams{
							Args:         []string{"api_group_1"},
							OperatorName: NotEqualOperator,
						},
					},
				},
			},
			expectedOutput: []string{"rule-2", "rule-3"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				Filters: []Filter{
					&APIGroupFilter{
						Params: FilterParams{
							Args:         []string{"api_group_1", "api_group_3", "api_group_2"},
							OperatorName: EqualsOperator,
						},
					},
				},
			},
			expectedOutput: []string{"rule-2"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				Filters: []Filter{
					&APIGroupFilter{
						Params: FilterParams{
							Args:         []string{"api_group_1", "api_group_3", "api_group_2"},
							OperatorName: EqualsOperator,
						},
					},
				},
			},
			expectedOutput: []string{"rule-2"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				Filters: []Filter{
					&NonResourceURLsFilter{
						Params: FilterParams{
							Args:         []string{},
							OperatorName: ExistsOperator,
						},
					},
				},
			},
			expectedOutput: []string{"rule-3"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				Filters: []Filter{
					&NonResourceURLsFilter{
						Params: FilterParams{
							Args:         []string{},
							OperatorName: DoesNotExistOperator,
						},
					},
				},
			},
			expectedOutput: []string{"rule-1", "rule-2"},
		},
		{
			input: RuleFilter{
				PermissionType: NameSpacedPermissionType,
				Filters: []Filter{
					&APIGroupFilter{
						Params: FilterParams{
							Args:         []string{"api_group_1"},
							OperatorName: InOperator,
						},
					},
				},
			},
			expectedOutput: []string{},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				Filters: []Filter{
					&APIGroupFilter{
						Params: FilterParams{
							Args:         []string{"api_group_1"},
							OperatorName: InOperator,
						},
					},
					&ResourcesFilter{
						Params: FilterParams{
							Args:         []string{"*"},
							OperatorName: InOperator,
						},
					},
				},
			},
			expectedOutput: []string{"rule-2"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				Filters: []Filter{
					&APIGroupFilter{
						Params: FilterParams{
							Args:         []string{"api_group_1"},
							OperatorName: NotEqualOperator,
						},
					},
					&VerbsFilter{
						Params: FilterParams{
							Args:         []string{"*"},
							OperatorName: InOperator,
						},
					},
				},
			},
			expectedOutput: []string{"rule-2", "rule-3"},
		},
		{
			input: RuleFilter{
				PermissionType: AllPermissionType,
				Filters: []Filter{
					&APIGroupFilter{
						Params: FilterParams{
							Args:         []string{"api_group_3", "api_group_1"},
							OperatorName: AnyOperator,
						},
					},
				},
			},
			expectedOutput: []string{"rule-1", "rule-2"},
		},
	}

	for _, testCase := range filterCases {
		tc := testCase
		t.Run("TestFilterRules",
			func(t *testing.T) {
				t.Parallel()
				res := inputRbacRules.FilterRules(tc.input)
				resRuleNames := make([]string, 0)
				for _, item := range res {
					resRuleNames = append(resRuleNames, item.name)
				}
				require.Equal(t, tc.expectedOutput, resRuleNames)
			})
	}
}
