package validators

import (
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
		{Code: "TEST0001"},
		{Code: "TEST0001"},
	}
	for _, v := range allValidators {
		registry.Add(v)
	}
}

func TestRegistryNonConcurrent(t *testing.T) {
	t.Parallel()
	registry := NewValidatorsRegistry()
	allValidators := []types.Validator{
		{Code: "TEST0001"},
		{Code: "TEST0002"},
		{Code: "TEST0003"},
		{Code: "TEST0004"},
		{Code: "TEST0005"},
		{Code: "TEST0006"},
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
