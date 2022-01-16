package validators

import (
	"errors"
	"testing"
	"time"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetryMiddlewareInterfaces(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(types.Middleware), RetryMiddleware{})
}

// Test the retry middleware, using a counter inside a closure
func TestRetryMiddleware(t *testing.T) {
	t.Parallel()

	expectedRetries := 5

	// begin count at -1 to ignore first attempt before entering retry loop
	retries := -1
	closure := types.ValidateFunc(func(types.ValidatorConfig, types.MetaBundle) types.ValidatorResult {
		retries++
		return RetryableError(errors.New("simply increase counter"))
	})

	v := types.NewValidator(
		"AM0000",
		closure,
		types.ValidatorName("Retry Counter"),
		types.ValidatorDescription("Count number of retries."),
		types.ValidatorRetryParams(&types.RetryParams{
			Count:   expectedRetries,
			Delayer: Constant(time.Duration(0)), // make retry faster
		}),
	)

	v.ApplyMiddleware(RetryMiddleware{})

	res := v.Validate(testutils.EmptyMetaBundle())
	require.True(t, res.IsRetryableError())
	require.Equal(t, expectedRetries, retries)
}

func TestConstantInterfaces(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(types.Delayer), new(Constant))
}

func TestConstantDelay(t *testing.T) {
	t.Parallel()

	delayer := Constant(2 * time.Second)

	require.Equal(t, delayer.Delay(), delayer.Delay())
}

func TestExponentialBackoffInterfaces(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(types.Delayer), new(ExponentialBackoff))
}

func TestExponentialBackoffDelay(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		delayer *ExponentialBackoff
		results []time.Duration
	}{
		"defauts": {
			delayer: NewExponentialBackoff(),
			results: []time.Duration{
				500 * time.Millisecond,
				1000 * time.Millisecond,
				2000 * time.Millisecond,
			},
		},
		"upper limit": {
			delayer: NewExponentialBackoff(
				ExponentialBackoffUpperLimit(499 * time.Millisecond),
			),
			results: []time.Duration{
				499 * time.Millisecond,
				499 * time.Millisecond,
				499 * time.Millisecond,
			},
		},
		"lower limit": {
			delayer: NewExponentialBackoff(
				ExponentialBackoffLowerLimit(501 * time.Millisecond),
			),
			results: []time.Duration{
				501 * time.Millisecond,
				1000 * time.Millisecond,
				2000 * time.Millisecond,
			},
		},
		"full jitter": {
			delayer: NewExponentialBackoff(
				ExponentialBackoffJitter(JitterFull),
				ExponentialBackoffRandomizer(MockDurationRandomizer(1*time.Second)),
			),
			results: []time.Duration{
				1 * time.Second,
				1 * time.Second,
				1 * time.Second,
			},
		},
		"equal jitter": {
			delayer: NewExponentialBackoff(
				ExponentialBackoffJitter(JitterEqual),
				ExponentialBackoffRandomizer(MockDurationRandomizer(1*time.Second)),
			),
			results: []time.Duration{
				1250 * time.Millisecond,
				1500 * time.Millisecond,
				2000 * time.Millisecond,
			},
		},
		"decorrellatted jitter": {
			delayer: NewExponentialBackoff(
				ExponentialBackoffJitter(JitterDecorrellated),
				ExponentialBackoffRandomizer(MockDurationRandomizer(1*time.Second)),
			),
			results: []time.Duration{
				500 * time.Millisecond,
				1500 * time.Millisecond,
				1500 * time.Millisecond,
			},
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			for _, res := range tc.results {
				assert.Equal(t, res, tc.delayer.Delay())
			}
		})
	}
}

type MockDurationRandomizer time.Duration

func (m MockDurationRandomizer) RandDurationBefore(time.Duration) time.Duration {
	return time.Duration(m)
}

func TestDefaultDurationRandomizerInterfaces(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(DurationRandomizer), new(DefaultDurationRandomizer))
}
