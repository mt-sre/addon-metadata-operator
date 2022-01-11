package validators

import (
	"time"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

func RetryMiddleware(v types.Validator) types.Validator {
	fn := func(mb types.MetaBundle) types.ValidatorResult {
		retryCount := 5
		retryDelaySeconds := 2

		if v.RetryCount > 0 {
			retryCount = v.RetryCount
		}
		if v.RetryDelaySeconds > 0 {
			retryDelaySeconds = v.RetryDelaySeconds
		}

		var res types.ValidatorResult
		for retryCount > 0 {
			res = v.Runner(mb)
			if !res.IsRetryableError() {
				return res
			}
			time.Sleep(time.Duration(retryDelaySeconds) * time.Second)
			retryCount--
		}
		// At this point retryCount == 0 and err is still retryable
		return res
	}

	return v.WithRunner(fn)
}

// CLIMiddlewares - wraps a validator with all CLI Middlewares. Order is important!
func CLIMiddlewares(v types.Validator) types.Validator {
	mws := []types.Middleware{
		RetryMiddleware,
	}
	return applyMiddlewares(mws, v)
}

func applyMiddlewares(mws []types.Middleware, v types.Validator) types.Validator {
	res := v
	for _, mw := range mws {
		res = mw(res)
	}
	return res
}
