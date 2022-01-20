package csvutils

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestWildCardApiGroupPresent(t *testing.T) {
	rule1 := types.CsvPermissions{
		ClusterPermissions: []types.Permissions{
			{
				Rules: []types.Rule{
					{
						ApiGroups:       []string{"api_group_1"},
						Resources:       []string{"resource_1"},
						Verbs:           []string{"*"},
						ResourceNames:   []string{"sample1"},
						NonResourceURLs: []string{},
					},
					{

						ApiGroups:       []string{"api_group_2", "api_group_3", "api_group_1"},
						Resources:       []string{"resource_2", "resource_3", "*"},
						Verbs:           []string{"*"},
						ResourceNames:   []string{"sample2", "sample3"},
						NonResourceURLs: []string{},
					},
				},
			},
		},
		Permissions: []types.Permissions{
			{
				Rules: []types.Rule{
					{
						ApiGroups:       []string{"*"},
						Resources:       []string{"resource_5"},
						Verbs:           []string{"*"},
						ResourceNames:   []string{"sample4"},
						NonResourceURLs: []string{"port-forward"},
					},
				},
			},
		},
	}
	rule2 := types.CsvPermissions{
		ClusterPermissions: []types.Permissions{
			{
				Rules: []types.Rule{
					{
						ApiGroups:       []string{"api_group_1"},
						Resources:       []string{"resource_1"},
						Verbs:           []string{"*"},
						ResourceNames:   []string{"sample1"},
						NonResourceURLs: []string{},
					},
					{

						ApiGroups:       []string{"api_group_2", "api_group_3", "api_group_1"},
						Resources:       []string{"resource_2", "resource_3", "*"},
						Verbs:           []string{"*"},
						ResourceNames:   []string{"sample2", "sample3"},
						NonResourceURLs: []string{},
					},
				},
			},
		},
		Permissions: []types.Permissions{
			{
				Rules: []types.Rule{
					{
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

	testCases := []struct {
		input          types.CsvPermissions
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
		res := WildCardApiGroupPresent(&testCase.input)
		require.Equal(t, testCase.expectedOutput, res)
	}
}

func TestWildCardResourcePresent(t *testing.T) {
	rule1 := types.CsvPermissions{
		ClusterPermissions: []types.Permissions{
			{
				Rules: []types.Rule{
					{
						ApiGroups:       []string{"api_group_1"},
						Resources:       []string{"*"},
						Verbs:           []string{"*"},
						ResourceNames:   []string{"sample1"},
						NonResourceURLs: []string{},
					},
				},
			},
		},
	}
	rule2 := types.CsvPermissions{
		ClusterPermissions: []types.Permissions{
			{
				Rules: []types.Rule{
					{

						ApiGroups:       []string{"api_group_2", "api_group_3", "api_group_1"},
						Resources:       []string{"*"},
						Verbs:           []string{"*"},
						ResourceNames:   []string{"sample2", "sample3"},
						NonResourceURLs: []string{},
					},
				},
			},
		},
	}

	rule3 := types.CsvPermissions{
		ClusterPermissions: []types.Permissions{
			{
				Rules: []types.Rule{
					{

						ApiGroups:       []string{""},
						Resources:       []string{"*"},
						Verbs:           []string{"*"},
						NonResourceURLs: []string{},
					},
				},
			},
		},
	}

	testCases := []struct {
		inputRule      types.CsvPermissions
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
		res := WildCardResourcePresent(&testCase.inputRule, testCase.ownedApis)
		require.Equal(t, testCase.expectedOutput, res)
	}
}

func TestCheckForConfidentialObjAccessAtClusterScope(t *testing.T) {
	rule1 := types.CsvPermissions{
		ClusterPermissions: []types.Permissions{
			{
				Rules: []types.Rule{
					{
						ApiGroups:       []string{""},
						Resources:       []string{"pods", "secrets", "configmaps"},
						Verbs:           []string{"*"},
						ResourceNames:   []string{"sample1"},
						NonResourceURLs: []string{},
					},
				},
			},
		},
	}
	rule2 := types.CsvPermissions{
		ClusterPermissions: []types.Permissions{
			{
				Rules: []types.Rule{
					{

						ApiGroups:       []string{""},
						Resources:       []string{"secrets"},
						Verbs:           []string{"read"},
						NonResourceURLs: []string{},
					},
				},
			},
		},
	}

	rule3 := types.CsvPermissions{
		ClusterPermissions: []types.Permissions{
			{
				Rules: []types.Rule{
					{

						ApiGroups:       []string{""},
						Resources:       []string{"configmaps"},
						Verbs:           []string{"update"},
						NonResourceURLs: []string{},
					},
				},
			},
		},
	}

	rule4 := types.CsvPermissions{
		ClusterPermissions: []types.Permissions{
			{
				Rules: []types.Rule{
					{

						ApiGroups:       []string{""},
						Resources:       []string{"configmaps", "secrets"},
						Verbs:           []string{"update"},
						NonResourceURLs: []string{},
					},
				},
			},
		},
	}

	rule5 := types.CsvPermissions{
		Permissions: []types.Permissions{
			{
				Rules: []types.Rule{
					{

						ApiGroups:       []string{""},
						Resources:       []string{"configmaps", "secrets"},
						Verbs:           []string{"update"},
						NonResourceURLs: []string{},
					},
				},
			},
		},
	}

	testCases := []struct {
		inputRule      types.CsvPermissions
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
		res := CheckForConfidentialObjAccessAtClusterScope(&testCase.inputRule)
		require.Equal(t, testCase.expectedOutput, res)
	}
}
