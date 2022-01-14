package validators

import (
	"context"

	amv1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
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
	SKURuleGetter
	// Stuck here until OCM Client can be injected via Params...
	CloseConnection() error
}

type SKURuleGetter interface {
	// GetSKURules returns any SKU Rules available in OCM which correspond
	// to the given OCM Quota name. If no SKU rules exist for a given OCM
	// quota name an empty slice is returned. An error is optionally
	// returned for any HTTP or network related issues are encountered.
	GetSKURules(context.Context, string) ([]*amv1.SkuRule, error)
}
