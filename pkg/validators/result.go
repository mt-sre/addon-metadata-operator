package validators

import "github.com/mt-sre/addon-metadata-operator/pkg/types"

func Success() types.ValidatorResult {
	return types.ValidatorResult{Success: true}
}

func Fail(failureMsg string) types.ValidatorResult {
	return types.ValidatorResult{Success: false, FailureMsg: failureMsg}
}

func Error(err error) types.ValidatorResult {
	return types.ValidatorResult{Success: false, Error: err, RetryableError: false}
}

// RetryableError - used by the Retry middleware re-run the validator
func RetryableError(err error) types.ValidatorResult {
	return types.ValidatorResult{Success: false, Error: err, RetryableError: true}
}
