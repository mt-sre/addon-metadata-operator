package validators

import (
	"errors"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/stretchr/testify/require"
)

var (
	emptyMetaBundle = types.MetaBundle{
		AddonMeta: &v1alpha1.AddonMetadataSpec{
			ID: "random-operator",
		},
	}
)

// Test the retry middleware, using a counter inside a closure
func TestRetryMiddleware(t *testing.T) {
	t.Parallel()
	retries := 0
	expectedRetries := 5
	closure := func(mb types.MetaBundle) types.ValidatorResult {
		retries++
		return RetryableError(errors.New("simply increase counter"))
	}
	v := types.Validator{
		Code:              "AM0000",
		Name:              "Retry Counter",
		Description:       "Count number of retries.",
		Runner:            closure,
		RetryCount:        expectedRetries,
		RetryDelaySeconds: 0, // make retry faster
	}
	v = RetryMiddleware(v)

	res := v.Runner(emptyMetaBundle)
	require.True(t, res.IsRetryableError())
	require.Equal(t, retries, expectedRetries)
}

// Test that all middlewares are appropriately registered using counter closure
func TestApplyMiddlewares(t *testing.T) {
	t.Parallel()
	counter := 0
	expectedCount := 5
	counterMiddleware := func(v types.Validator) types.Validator {
		fn := func(mb types.MetaBundle) types.ValidatorResult {
			counter++
			return v.Runner(mb)
		}
		return v.WithRunner(fn)
	}
	v := types.Validator{
		Code:        "AM0000",
		Name:        "Success Validator",
		Description: "Simply succeeds.",
		Runner:      func(mb types.MetaBundle) types.ValidatorResult { return Success() },
	}

	mws := make([]types.Middleware, 0)
	for i := 0; i < expectedCount; i++ {
		mws = append(mws, counterMiddleware)
	}
	v = applyMiddlewares(mws, v)

	res := v.Runner(emptyMetaBundle)
	require.True(t, res.IsSuccess())
	require.Equal(t, counter, expectedCount)
}
