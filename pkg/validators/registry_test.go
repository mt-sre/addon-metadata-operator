package validators

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestRegistryPanicsDuplicateKey(t *testing.T) {
	t.Parallel()
	defer func() {
		require.NotNil(t, recover())
	}()

	registry := NewValidatorsRegistry()
	allValidators := []types.Validator{
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM0001",
			},
		},
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM0001",
			},
		},
	}
	for _, v := range allValidators {
		registry.Add(v)
	}
}

func TestRegistryNonConcurrent(t *testing.T) {
	t.Parallel()
	registry := NewValidatorsRegistry()
	allValidators := []types.Validator{
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM0001",
			},
		},
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM0002",
			},
		},
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM00003",
			},
		},
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM00004",
			},
		},
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM00005",
			},
		},
	}
	require.Equal(t, registry.Len(), 0)
	for _, v := range allValidators {
		registry.Add(v)
		vCopy, ok := registry.Get(v.Code)
		require.True(t, ok)
		require.Equal(t, v, vCopy)
	}
	require.Equal(t, registry.Len(), len(allValidators))
}

func TestRegistryListSorted(t *testing.T) {
	t.Parallel()
	const numShuffles = 10
	allValidators := []types.Validator{
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM0001",
			},
		},
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM0000",
			},
		},
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM0005",
			},
		},
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM0004",
			},
		},
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM0003",
			},
		},
		{
			ValidatorConfig: types.ValidatorConfig{
				Code: "AM0002",
			},
		},
	}

	for i := 0; i < numShuffles; i++ {
		rand.Shuffle(len(allValidators), func(i, j int) {
			allValidators[i], allValidators[j] = allValidators[j], allValidators[i]
		})

		ensureValidatorOrder(t, allValidators)
	}
}

func ensureValidatorOrder(t *testing.T, allValidators []types.Validator) {
	t.Helper()

	registry := NewValidatorsRegistry()

	for _, v := range allValidators {
		registry.Add(v)
	}

	for i, v := range registry.ListSorted() {
		require.Equal(t, fmt.Sprintf("AM%04d", i), v.Code)
	}
}
