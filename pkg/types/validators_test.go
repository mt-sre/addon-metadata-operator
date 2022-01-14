package types_test

import (
	"sort"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestValidatorListImplementsSort(t *testing.T) {
	require.Implements(t, (*sort.Interface)(nil), types.ValidatorList{})
}
