package validators

import (
	"context"
)

// OCMError abstracts behavior required for validators to identify the underlying
// causes of OCM related errors.
type OCMError interface {
	// ServerSide returns 'true' if an instance of OCMError was caused
	// by a server-side issue.
	ServerSide() bool
}

// IsOCMServerSideError determines if the given error is both an instance of OCMError
// and was caused by a server-side issue.
func IsOCMServerSideError(err error) bool {
	ocmErr, ok := err.(OCMError)

	return ok && ocmErr.ServerSide()
}

// OCMClient abstracts behavior required for validators which request data
// from OCM to be implemented by OCM API clients.
type OCMClient interface {
	QuotaRuleGetter
	// Stuck here until OCM Client can be injected via Params...
	CloseConnection() error
}

type QuotaRuleGetter interface {
	QuotaRuleExists(context.Context, string) (bool, error)
}
