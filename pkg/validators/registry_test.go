package validators

import (
	"math/rand"
	"strconv"
	"strings"
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

func TestRegistryListSorted(t *testing.T) {
	t.Parallel()
	const numShuffles = 10
	allValidators := []types.Validator{
		{Code: "TEST0001"},
		{Code: "TEST0000"},
		{Code: "TEST0005"},
		{Code: "TEST0004"},
		{Code: "TEST0003"},
		{Code: "TEST0002"},
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
	allValidatorsSorted := registry.ListSorted()
	for i := 0; i < len(allValidatorsSorted); i++ {
		code := allValidatorsSorted[i].Code
		parts := strings.Split(code, "TEST")
		require.Equal(t, len(parts), 2)

		j, err := strconv.Atoi(parts[1])
		require.NoError(t, err)
		require.Equal(t, i, j)
	}
}
