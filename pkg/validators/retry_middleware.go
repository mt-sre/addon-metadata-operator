package validators

import (
	"math/rand"
	"time"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

// RetryMiddleware implements middleware that conditionally
// retries types.ValidatorRunner executions based on the result.
type RetryMiddleware struct{}

func (r RetryMiddleware) Wrap(v types.ValidatorRunner) types.ValidatorRunner {
	return types.ValidateFunc(func(cfg types.ValidatorConfig, mb types.MetaBundle) types.ValidatorResult {
		var (
			maxRetries               = 5
			delayer    types.Delayer = Constant(2 * time.Second)
		)

		if cfg.RetryParams != nil {
			maxRetries = cfg.RetryParams.Count
			delayer = cfg.RetryParams.Delayer
		}

		res := v.Run(cfg, mb)

		for i := 0; i < maxRetries; i++ {
			if !res.IsRetryableError() {
				return res
			}

			time.Sleep(delayer.Delay())

			res = v.Run(cfg, mb)
		}

		return res
	})
}

// Constant implements a constant-time delayer which will
// always return the value from which it is instantiated.
type Constant time.Duration

func (d Constant) Delay() time.Duration {
	return time.Duration(d)
}

// NewExponentialBackoff applies any given options to an instance of
// ExponentialBackoff and provides defaults for those not provided.
// A *ExponentialBackoff instance is then returned ready to be used
// as a Delayer returning delays values corresponding to it's configuration.
func NewExponentialBackoff(opts ...ExponentialBackoffOption) *ExponentialBackoff {
	var delayer ExponentialBackoff

	for _, opt := range opts {
		opt(&delayer)
	}

	if delayer.base == 0 {
		delayer.base = 500 * time.Millisecond
	}

	if delayer.upperLimit == 0 {
		delayer.upperLimit = 5 * time.Second
	}

	delayer.current = durationMin(delayer.base, delayer.upperLimit)

	if delayer.rand == nil {
		delayer.rand = NewDefaultDurationRandomizer()
	}

	return &delayer
}

// ExponentialBackoff implements a types.Delayer which can be configured
// to successively generate delay values which either increase exponentially
// or are uniformally generated from a range derived from an exponentiated value.
type ExponentialBackoff struct {
	current    time.Duration
	base       time.Duration
	lowerLimit time.Duration
	upperLimit time.Duration
	jitter     Jitter
	rand       DurationRandomizer
}

func (d *ExponentialBackoff) Delay() time.Duration {
	defer d.increaseDelay()

	var delay time.Duration

	switch d.jitter {
	case JitterFull:
		delay = d.rand.RandDurationBefore(d.current)
	case JitterEqual:
		delay = d.current/2 + d.rand.RandDurationBefore(d.current/2)
	default:
		delay = d.current
	}

	return durationMax(d.lowerLimit, delay)
}

func (d *ExponentialBackoff) increaseDelay() {
	if d.current >= d.upperLimit {
		return
	}

	var nextValue time.Duration

	switch d.jitter {
	case JitterDecorrellated:
		nextValue = d.base + d.rand.RandDurationBefore(d.current*3)
	default:
		nextValue = d.current << 1
	}

	d.current = durationMin(d.upperLimit, nextValue)
}

// ExponentialBackoffOption abstracts an option to apply to an
// instance of ExponentialBackoff.
type ExponentialBackoffOption func(d *ExponentialBackoff)

// ExponentialBackoffBaseDelay sets the base delay value from which
// delay values will start generating.
func ExponentialBackoffBaseDelay(base time.Duration) ExponentialBackoffOption {
	return func(d *ExponentialBackoff) {
		d.base = base
	}
}

// ExponentialBackoffUpperLimit set the upper delay limit ensuring that delays
// do not increase infinitely.
func ExponentialBackoffUpperLimit(upper time.Duration) ExponentialBackoffOption {
	return func(d *ExponentialBackoff) {
		d.upperLimit = upper
	}
}

// ExponentialBackoffLowerLimit set the lower delay limit ensuring that delays
// do not fall below a minimum duration.
func ExponentialBackoffLowerLimit(lower time.Duration) ExponentialBackoffOption {
	return func(d *ExponentialBackoff) {
		d.lowerLimit = lower
	}
}

// ExponentialBackoffRandomizer sets the DurationRandomizer instance which is
// responsible for generating randomness when applying jitter.
func ExponentialBackoffRandomizer(rand DurationRandomizer) ExponentialBackoffOption {
	return func(d *ExponentialBackoff) {
		d.rand = rand
	}
}

// ExponentialBackoffJitter sets the type of jitter to apply when generating
// new delay values.
func ExponentialBackoffJitter(jitter Jitter) ExponentialBackoffOption {
	return func(d *ExponentialBackoff) {
		d.jitter = jitter
	}
}

// DurationRandomizer abstracts behavior which generates random
// time.Duration values.
type DurationRandomizer interface {
	// RandDurationBefore returns a random time.Duration guaranteed
	// to be less than or equal to the given maximum.
	RandDurationBefore(max time.Duration) time.Duration
}

// NewDefaultDurationRandomizer returns an initalized instance of
// *DefaultDurationRandomizer which seeds it's randomness from
// the time as of initialization.
func NewDefaultDurationRandomizer() *DefaultDurationRandomizer {
	return &DefaultDurationRandomizer{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// DefaultDurationRandomizer provides random duration values
// generated from a discrete uniform distribution.
type DefaultDurationRandomizer struct {
	*rand.Rand
}

func (dr *DefaultDurationRandomizer) RandDurationBefore(max time.Duration) time.Duration {
	return time.Duration(dr.Int63n(max.Nanoseconds()))
}

// Jitter is an enum type describing jitter algorithms implemented in this package.
type Jitter int

// See https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter for more info
const (
	// JitterNone applies no randomness to exponential backoff values
	JitterNone Jitter = iota
	// JitterFull applies total randomness to exponential backoff values
	JitterFull
	// JitterEqual applies an equal proportion of randomness and determinism
	// to exponential backoff values. This appears to be a special case of
	// a more general class of jitter applications which "tunes" the ratio
	// of jitter.
	JitterEqual
	// JitterDecorrellated is not based in exponentiation, but rather scales with
	// prior random results ala Markov Chains. The choice of '3' as the constant
	// multiplier in this algorithm is unexplained, but it likely was chosen to
	// achieve an average curvature more similiar to 2^n.
	JitterDecorrellated
)

func durationMin(d1, d2 time.Duration) time.Duration {
	if d1 < d2 {
		return d1
	}

	return d2
}

func durationMax(d1, d2 time.Duration) time.Duration {
	if d1 >= d2 {
		return d1
	}

	return d2
}
