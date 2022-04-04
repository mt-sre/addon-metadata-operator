package validator

import (
	"context"
	"testing"
	"time"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCode(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		Input          string
		Expected       Code
		ErrorAssertion assert.ErrorAssertionFunc
	}{
		"one": {
			Input:          "AM0001",
			Expected:       Code(1),
			ErrorAssertion: assert.NoError,
		},
		"one-thousand": {
			Input:          "AM1000",
			Expected:       Code(1000),
			ErrorAssertion: assert.NoError,
		},
		"lower case prefix": {
			Input:          "am0001",
			Expected:       Code(0001),
			ErrorAssertion: assert.NoError,
		},
		"more than 6 characters": {
			Input:          "AM10000",
			Expected:       Code(0),
			ErrorAssertion: assert.Error,
		},
		"less than four zero padding": {
			Input:          "AM001",
			Expected:       Code(0),
			ErrorAssertion: assert.Error,
		},
		"wrong prefix": {
			Input:          "PM1000",
			Expected:       Code(0),
			ErrorAssertion: assert.Error,
		},
		"arbitrary 6 character string": {
			Input:          "abcdef",
			Expected:       Code(0),
			ErrorAssertion: assert.Error,
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			code, err := ParseCode(tc.Input)
			tc.ErrorAssertion(t, err)

			assert.Equal(t, tc.Expected, code)
		})
	}
}

func TestRunnerRegistration(t *testing.T) {
	t.Parallel()

	const (
		code = Code(0)
		name = "dummy_validator"
		desc = "this is a dummy validator"
	)

	Register(
		NewValidatorMock(
			code,
			name,
			desc,
			func(context.Context, types.MetaBundle) Result {
				return Result{success: true}
			},
		),
	)

	runner, err := NewRunner(
		WithOCMClient{testutils.NewMockOCMClient()},
	)
	require.NoError(t, err)

	vals := runner.GetValidators(MatchesCodes(code))
	assert.Len(t, vals, 1)
	assert.Equal(t, code, vals[0].Code())
	assert.Equal(t, name, vals[0].Name())
	assert.Equal(t, desc, vals[0].Description())

	vals = runner.GetValidators(Not(MatchesCodes(code)))
	assert.Len(t, vals, 0)
}

func TestRunnerMiddleware(t *testing.T) {
	t.Parallel()

	const (
		code = Code(0)
		name = "dummy_validator"
		desc = "this is a dummy validator"
	)

	const expectedCount = 3

	var actualCount int

	runner, err := NewRunner(
		WithOCMClient{testutils.NewMockOCMClient()},
		WithInitializers{
			NewValidatorMock(
				code,
				name,
				desc,
				func(context.Context, types.MetaBundle) Result {
					actualCount++

					return Result{retryable: true}
				},
			),
		},
		WithMiddleware{
			NewRetryMiddleware(
				WithMaxAttempts(3),
				WithDelay(0*time.Second),
			),
		},
	)
	require.NoError(t, err)

	<-runner.Run(context.TODO(), types.MetaBundle{})
	assert.Equal(t, expectedCount, actualCount)
}

func NewValidatorMock(
	code Code,
	name, desc string,
	runner func(context.Context, types.MetaBundle) Result) func(Dependencies) (Validator, error) {

	base, err := NewBase(
		code,
		BaseName(name),
		BaseDesc(desc),
	)

	return func(Dependencies) (Validator, error) {
		return &ValidatorMock{
			Base:   base,
			runner: runner,
		}, err
	}
}

type ValidatorMock struct {
	*Base
	runner func(context.Context, types.MetaBundle) Result
}

func (v ValidatorMock) Run(ctx context.Context, mb types.MetaBundle) Result {
	return v.runner(ctx, mb)
}
