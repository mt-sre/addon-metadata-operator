package validate

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Validator interface {
	Name() string
	Run(utils.MetaBundle) (bool, error)
	SucceedingCandidates() []utils.MetaBundle
	FailingCandidates() []utils.MetaBundle
}

// register the validators to test here by appending them to the following slice
var validatorsToTest []Validator = []Validator{
	validators.ValidatorAddonLabelTestBundle{},
	validators.ValidatorDefaultChannelTestBundle{},
}

func Test_AllValidators(t *testing.T) {
	for _, validator := range validatorsToTest {
		validator := validator
		t.Run(validator.Name(), func(t *testing.T) {
			t.Parallel()
			// testing the succeeding candidates
			succeedingMetaBundles := validator.SucceedingCandidates()
			for _, mb := range succeedingMetaBundles {
				success, err := validator.Run(mb)
				require.NoError(t, err)
				assert.True(t, success)
			}

			// testing the failing candidates
			failingMetaBundles := validator.FailingCandidates()
			for _, mb := range failingMetaBundles {
				success, _ := validator.Run(mb)
				assert.False(t, success)
			}
		})
	}
}
