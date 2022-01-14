package validators

import "github.com/mt-sre/addon-metadata-operator/pkg/types"

func Success() types.ValidatorResult {
	return types.ValidatorResultSuccess()
}

func Fail(msgs ...string) types.ValidatorResult {
	return types.ValidatorResultFailure(msgs...)
}

func Error(err error) types.ValidatorResult {
	return types.ValidatorResultError(err, false)
}

// RetryableError - used by the Retry middleware to automatically re-run validators
func RetryableError(err error) types.ValidatorResult {
	return types.ValidatorResultError(err, true)
}
