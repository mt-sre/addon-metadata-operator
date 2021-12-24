package validate

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// All validators implement the ValidatorTest interface
// register the validators to test here by appending them to the following slice
var validatorsToTest []utils.ValidatorTest = []utils.ValidatorTest{
	validators.ValidatorTest001DefaultChannel{},
	validators.ValidatorAddonLabelTestBundle{},
}

func TestAllValidators(t *testing.T) {
	for _, validator := range validatorsToTest {
		validator := validator
		t.Run(validator.Name(), func(t *testing.T) {
			t.Parallel()
			// testing the succeeding candidates
			succeedingMetaBundles := validator.SucceedingCandidates()
			for _, mb := range succeedingMetaBundles {
				success, failureMsg, err := validator.Run(mb)
				require.NoError(t, err)
				assert.True(t, success, failureMsg)
			}

			// (optional) testing the failing candidates
			failingMetaBundles := validator.FailingCandidates()
			for _, mb := range failingMetaBundles {
				success, failureMsg, _ := validator.Run(mb)
				assert.False(t, success, failureMsg)
			}
		})
	}
}

func TestFilterDisabledValidators(t *testing.T) {
	n_validators := len(AllValidators)

	cases := []struct {
		name     string
		disabled []string
	}{
		{
			name:     "all_enabled",
			disabled: []string{},
		},
		{
			name:     "disable_default_channel",
			disabled: []string{"001_default_channel"},
		},
		{
			name:     "disable_two",
			disabled: []string{"001_default_channel", "002_label_format"},
		},
	}
	for _, tc := range cases {
		tc := tc // pin
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			filter, err := NewFilter(strings.Join(tc.disabled, ","), "")
			require.NoError(t, err)

			n_enabled := len(filter.GetValidators())
			n_disabled := len(tc.disabled)
			require.Equal(t, n_enabled+n_disabled, n_validators)
		})
	}
}

func TestFilterEnabledValidators(t *testing.T) {
	cases := []struct {
		name    string
		enabled []string
	}{
		{
			name:    "enable_default_channel",
			enabled: []string{"001_default_channel"},
		},
		{
			name:    "enable_two",
			enabled: []string{"001_default_channel", "002_label_format"},
		},
	}
	for _, tc := range cases {
		tc := tc // pin
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			filter, err := NewFilter("", strings.Join(tc.enabled, ","))
			require.NoError(t, err)
			require.Equal(t, len(filter.GetValidators()), len(tc.enabled))
		})
	}
}

func TestEmptyFilterAllEnabled(t *testing.T) {
	t.Parallel()
	filter, err := NewFilter("", "")
	require.NoError(t, err)
	require.Equal(t, len(filter.GetValidators()), len(AllValidators))
}

func TestFilterError(t *testing.T) {
	cases := []struct {
		name     string
		enabled  []string
		disabled []string
	}{
		{
			name:     "mutually_exclusive",
			enabled:  []string{"001_default_channel"},
			disabled: []string{"001_default_channel"},
		},
		{
			name:     "enabled_dont_exist",
			enabled:  []string{"invalid"},
			disabled: []string{},
		},
		{
			name:     "disabled_dont_exist",
			enabled:  []string{},
			disabled: []string{"invalid"},
		},
	}
	for _, tc := range cases {
		tc := tc // pin
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			disabled := strings.Join(tc.disabled, ",")
			enabled := strings.Join(tc.enabled, ",")
			filter, err := NewFilter(disabled, enabled)
			require.Error(t, err)
			require.Nil(t, filter)
		})
	}
}

func TestAllValidatorsNamesAreUnique(t *testing.T) {
	t.Parallel()
	seen := make(map[string]int)
	for name, validator := range AllValidators {
		require.Equal(t, name, validator.Name, "Name %v and %v don't match in AllValidators map.", name, validator.Name)
		seen[name]++
	}
	for name, count := range seen {
		require.Equal(t, count, 1, fmt.Sprintf("Validator name %v is not unique.", name))
	}
}

// TODO (sblaisdo) - enable after validators development stabilizes
// little math trick to make sure we have an arithmetic sequence of n terms (Gauss)
// func TestAllValidatorsNamesFollowArithmeticSequence(t *testing.T) {
// 	n := len(AllValidators)
// 	sum := 0
// 	for name := range AllValidators {
// 		parts := strings.SplitN(name, "_", 2)
// 		i, err := strconv.Atoi(parts[0])
// 		require.NoError(t, err)
// 		sum += i
// 	}
// 	require.Equal(t, sum, (n*n+n)/2, "Please make sure validator names follow an arithmetic sequence.")
// }
