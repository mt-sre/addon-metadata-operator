package validator

// Result encapsulates the status and reason for the result of
// a Validator task running against a types.MetaBundle.
type Result struct {
	Code        Code
	Name        string
	Description string
	FailureMsgs []string
	Error       error
	retryable   bool
	success     bool
}

// IsSuccess returns 'true' if the Validator task which
// returned it was successful.
func (r Result) IsSuccess() bool { return r.success }

// IsError returns 'true' if the Validator task which
// returned it encountered an error.
func (r Result) IsError() bool { return r.Error != nil }

// IsRetryableError returns 'true' if the Validator task which
// returned it encountered an error, but the error can be retried.
func (r Result) IsRetryableError() bool { return r.retryable }

// ResultList is a sortable slice of Result instances.
type ResultList []Result

func (l ResultList) Len() int           { return len(l) }
func (l ResultList) Less(i, j int) bool { return l[i].Code < l[j].Code }
func (l ResultList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

// HasFailure returns 'true' if any of the ResultList members
// are failures or errors.
func (l ResultList) HasFailure() bool {
	for _, r := range l {
		if r.IsSuccess() {
			continue
		}

		return true
	}

	return false
}

// Errors returns a slice of errors from the ResultList
// members. If no errors were encountered then an empty slice
// is returned.
func (l ResultList) Errors() []error {
	var errs []error

	for _, r := range l {
		if !r.IsError() {
			continue
		}

		errs = append(errs, r.Error)
	}

	return errs
}
