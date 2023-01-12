package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOCMClientImplInterface(t *testing.T) {
	t.Parallel()

	require.Implements(t,
		new(OCMClient), new(OCMClientImpl),
		"utils.OCMClientImpl must implement the types.OCMClient interface",
	)
}

func TestOCMResponseErrorInterface(t *testing.T) {
	t.Parallel()

	require.Implements(t,
		new(OCMError), OCMResponseError(400),
		"utils.OCMResponseError must implement the validator.OCMError interface",
	)
}

func TestOCMResponseErrorServerSide(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		code   int
		result bool
	}{
		"negative one": {
			code: -1, result: false,
		},
		"zero": {
			code: 0, result: false,
		},
		"out of range": {
			code: 600, result: false,
		},
		"OK": {
			code: 200, result: false,
		},
		"MOVED PERMANENTLY": {
			code: 301, result: false,
		},
		"BAD REQUEST": {
			code: 400, result: false,
		},
		"NOT FOUND": {
			code: 404, result: false,
		},
		"INTERNAL SERVER ERROR": {
			code: 500, result: true,
		},
		"BAD GATEWAY": {
			code: 502, result: true,
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, OCMResponseError(tc.code).ServerSide(), tc.result)
		})
	}
}

func TestDisconnectedOCMClientInterface(t *testing.T) {
	t.Parallel()

	require.Implements(t,
		new(OCMClient), new(DisconnectedOCMClient),
		"utils.DisconnectedOCMClient must implement the types.OCMClient interface",
	)
}

func TestDisconnectedOCMClientQuotaRuleExists(t *testing.T) {
	t.Parallel()

	var client DisconnectedOCMClient

	_, err := client.QuotaRuleExists(context.Background(), "")
	require.Error(t, err)
}
