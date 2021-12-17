package validate

import (
	"strings"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Validator interface {
	Name() string
	Run(utils.MetaBundle) (bool, string, error)
	SucceedingCandidates() []utils.MetaBundle
	FailingCandidates() []utils.MetaBundle
}

// register the validators to test here by appending them to the following slice
var validatorsToTest []Validator = []Validator{
	validators.Validator001DefaultChannel{},
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

			// testing the failing candidates
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
			name:     "disable_all",
			disabled: []string{"001_default_channel", "002_label_format", "003_csv_present"},
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
