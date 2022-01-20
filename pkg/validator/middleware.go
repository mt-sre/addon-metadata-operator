package validator

import (
	"context"
	"time"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

type Middleware interface {
	Wrap(RunFunc) RunFunc
}

type RunFunc func(context.Context, types.MetaBundle) Result

func NewRetryMiddleware(opts ...RetryMiddlewareOption) *RetryMiddleware {
	cfg := RetryMiddlewareConfig{
		Delay: 2 * time.Second,
	}

	cfg.Option(opts...)
	cfg.Default()

	return &RetryMiddleware{
		cfg: cfg,
	}
}

type RetryMiddleware struct {
	cfg RetryMiddlewareConfig
}

func (r *RetryMiddleware) Wrap(run RunFunc) RunFunc {
	return func(ctx context.Context, mb types.MetaBundle) Result {
		var res Result

		for attempts := 0; attempts < r.cfg.MaxAttempts; attempts++ {
			res = run(ctx, mb)
			if !res.IsRetryableError() {
				return res
			}

			time.Sleep(time.Duration(r.cfg.Delay))
		}

		return res
	}
}

type RetryMiddlewareConfig struct {
	MaxAttempts int
	Delay       time.Duration
}

func (c *RetryMiddlewareConfig) Option(opts ...RetryMiddlewareOption) {
	for _, opt := range opts {
		opt.ConfigureRetryMiddleware(c)
	}
}

func (c *RetryMiddlewareConfig) Default() {
	if c.MaxAttempts == 0 {
		c.MaxAttempts = 5
	}
}

type RetryMiddlewareOption interface {
	ConfigureRetryMiddleware(*RetryMiddlewareConfig)
}

type WithMaxAttempts int

func (ma WithMaxAttempts) ConfigureRetryMiddleware(c *RetryMiddlewareConfig) {
	c.MaxAttempts = int(ma)
}

type WithDelay time.Duration

func (d WithDelay) ConfigureRetryMiddleware(c *RetryMiddlewareConfig) {
	c.Delay = time.Duration(d)
}
