package csvutils

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/stretchr/testify/require"
	rbac "k8s.io/api/rbac/v1"
)

func TestWildCardApiGroupPresent(t *testing.T) {
	rule1 := types.CSVPermissions{
		ClusterPermissions: []types.Permission{
			{
				Rules: []types.Rule{
					{
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{"api_group_1"},
							Resources:       []string{"resource_1"},
							Verbs:           []string{"*"},
							ResourceNames:   []string{"sample1"},
							NonResourceURLs: []string{},
						},
					},
					{
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
		Permissions: []types.Permission{
			{
				Rules: []types.Rule{
					{
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{"*"},
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
	rule2 := types.CSVPermissions{
		ClusterPermissions: []types.Permission{
			{
				Rules: []types.Rule{
					{
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{"api_group_1"},
							Resources:       []string{"resource_1"},
							Verbs:           []string{"*"},
							ResourceNames:   []string{"sample1"},
							NonResourceURLs: []string{},
						},
					},
					{
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
		Permissions: []types.Permission{
			{
				Rules: []types.Rule{
					{
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

	testCases := []struct {
		input          types.CSVPermissions
		expectedOutput bool
	}{
		{
			input:          rule1,
			expectedOutput: true,
		},
		{
			input:          rule2,
			expectedOutput: false,
		},
	}

	for _, testCase := range testCases {
		tc := testCase // pin
		t.Run(
			"TestWildCardApiGroupPresent",
			func(t *testing.T) {
				t.Parallel()
				res := WildCardApiGroupPresent(&tc.input)
				require.Equal(t, tc.expectedOutput, res)
			},
		)
	}
}

func TestWildCardResourcePresent(t *testing.T) {
	rule1 := types.CSVPermissions{
		ClusterPermissions: []types.Permission{
			{
				Rules: []types.Rule{
					{
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{"api_group_1"},
							Resources:       []string{"*"},
							Verbs:           []string{"*"},
							ResourceNames:   []string{"sample1"},
							NonResourceURLs: []string{},
						},
					},
				},
			},
		},
	}
	rule2 := types.CSVPermissions{
		ClusterPermissions: []types.Permission{
			{
				Rules: []types.Rule{
					{
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{"api_group_2", "api_group_3", "api_group_1"},
							Resources:       []string{"*"},
							Verbs:           []string{"*"},
							ResourceNames:   []string{"sample2", "sample3"},
							NonResourceURLs: []string{},
						},
					},
				},
			},
		},
	}

	rule3 := types.CSVPermissions{
		ClusterPermissions: []types.Permission{
			{
				Rules: []types.Rule{
					{
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{""},
							Resources:       []string{"*"},
							Verbs:           []string{"*"},
							NonResourceURLs: []string{},
						},
					},
				},
			},
		},
	}

	testCases := []struct {
		inputRule      types.CSVPermissions
		ownedApis      []string
		expectedOutput bool
	}{
		{
			inputRule:      rule1,
			ownedApis:      []string{"api_group_1"},
			expectedOutput: false,
		},
		{
			inputRule:      rule2,
			ownedApis:      []string{"api_group_2", "api_group_3"},
			expectedOutput: true,
		},
		{
			inputRule:      rule3,
			ownedApis:      []string{"api_group_3"},
			expectedOutput: true,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(
			"TestWildCardResourcePresent",
			func(t *testing.T) {
				t.Parallel()
				res := WildCardResourcePresent(&tc.inputRule, tc.ownedApis)
				require.Equal(t, tc.expectedOutput, res)
			},
		)
	}
}

func TestCheckForConfidentialObjAccessAtClusterScope(t *testing.T) {
	rule1 := types.CSVPermissions{
		ClusterPermissions: []types.Permission{
			{
				Rules: []types.Rule{
					{
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{""},
							Resources:       []string{"pods", "secrets", "configmaps"},
							Verbs:           []string{"*"},
							ResourceNames:   []string{"sample1"},
							NonResourceURLs: []string{},
						},
					},
				},
			},
		},
	}
	rule2 := types.CSVPermissions{
		ClusterPermissions: []types.Permission{
			{
				Rules: []types.Rule{
					{
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{""},
							Resources:       []string{"secrets"},
							Verbs:           []string{"read"},
							NonResourceURLs: []string{},
						},
					},
				},
			},
		},
	}

	rule3 := types.CSVPermissions{
		ClusterPermissions: []types.Permission{
			{
				Rules: []types.Rule{
					{
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{""},
							Resources:       []string{"configmaps"},
							Verbs:           []string{"update"},
							NonResourceURLs: []string{},
						},
					},
				},
			},
		},
	}

	rule4 := types.CSVPermissions{
		ClusterPermissions: []types.Permission{
			{
				Rules: []types.Rule{
					{
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{""},
							Resources:       []string{"configmaps", "secrets"},
							Verbs:           []string{"update"},
							NonResourceURLs: []string{},
						},
					},
				},
			},
		},
	}

	rule5 := types.CSVPermissions{
		Permissions: []types.Permission{
			{
				Rules: []types.Rule{
					{
						PolicyRule: rbac.PolicyRule{
							APIGroups:       []string{""},
							Resources:       []string{"configmaps", "secrets"},
							Verbs:           []string{"update"},
							NonResourceURLs: []string{},
						},
					},
				},
			},
		},
	}

	testCases := []struct {
		inputRule      types.CSVPermissions
		expectedOutput bool
	}{
		{
			inputRule:      rule1,
			expectedOutput: false,
		},
		{
			inputRule:      rule2,
			expectedOutput: true,
		},
		{
			inputRule:      rule3,
			expectedOutput: true,
		},
		{
			inputRule:      rule4,
			expectedOutput: true,
		},
		{
			inputRule:      rule5,
			expectedOutput: false,
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(
			"TestCheckForConfidentialObjAccessAtClusterScope",
			func(t *testing.T) {
				t.Parallel()
				res := CheckForConfidentialObjAccessAtClusterScope(&tc.inputRule)
				require.Equal(t, tc.expectedOutput, res)
			},
		)
	}
}
