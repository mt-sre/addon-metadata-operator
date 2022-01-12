package validators_test

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRegistry - register all test structs
var TestRegistry = NewTestRegistry()

func NewTestRegistry() *testRegistry {
	return &testRegistry{[]types.ValidatorTest{}}
}

type testRegistry struct {
	Data []types.ValidatorTest
}

// Add - update the test registry in a thread-safe way. Called in init() functions
func (t *testRegistry) Add(v types.ValidatorTest) {
	t.Data = append(t.Data, v)
}

func (t *testRegistry) All() []types.ValidatorTest {
	return t.Data
}
func TestAllValidators(t *testing.T) {
	for _, validator := range TestRegistry.All() {
		validator := validator
		t.Run(validator.Name(), func(t *testing.T) {
			t.Parallel()
			// testing the succeeding candidates
			succeedingMetaBundles := validator.SucceedingCandidates()
			for _, mb := range succeedingMetaBundles {
				res := validator.Run(mb)
				require.False(t, res.IsError())
				assert.True(t, res.IsSuccess())
			}

			// (optional) testing the failing candidates
			failingMetaBundles := validator.FailingCandidates()
			for _, mb := range failingMetaBundles {
				res := validator.Run(mb)
				require.False(t, res.IsError())
				assert.False(t, res.IsSuccess())
			}
		})
	}
}
