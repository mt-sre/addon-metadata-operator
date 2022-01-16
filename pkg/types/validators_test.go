package types

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test that all middlewares are appropriately registered using counter closure
func TestValidatorApplyMiddleware(t *testing.T) {
	t.Parallel()

	expectedCount := 5

	v := NewValidator(
		"AM0000",
		ValidateFunc(func(ValidatorConfig, MetaBundle) ValidatorResult {
			return ValidatorResult{Success: true}
		}),
		ValidatorName("Success Validator"),
		ValidatorDescription("Simply succeeds."),
	)

	counter := counterMiddleware(0)

	mws := make([]Middleware, 0, expectedCount)

	for i := 0; i < expectedCount; i++ {
		mws = append(mws, &counter)
	}

	v.ApplyMiddleware(mws...)

	res := v.Validate(MetaBundle{})
	require.True(t, res.IsSuccess())
	require.Equal(t, expectedCount, int(counter))
}

type counterMiddleware int

func (mw *counterMiddleware) Wrap(v ValidatorRunner) ValidatorRunner {
	return ValidateFunc(func(cfg ValidatorConfig, mb MetaBundle) ValidatorResult {
		*mw++

		return v.Run(cfg, mb)
	})
}

func TestValidateFuncInterfaces(t *testing.T) {
	t.Parallel()

	require.Implements(t,
		new(ValidatorRunner),
		ValidateFunc(func(ValidatorConfig, MetaBundle) ValidatorResult {
			return ValidatorResult{}
		}),
	)
}

func TestValidatorListImplementsSort(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(sort.Interface), ValidatorList{})
}
