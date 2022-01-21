package utils_test

import (
	"os"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultOCMClientInterface(t *testing.T) {
	t.Parallel()

	require.Implements(t,
		new(validators.OCMClient), new(utils.DefaultOCMClient),
		"utils.DefaultOCMClient must implement the types.OCMClient interface",
	)
}

func TestOCMResponseErrorInterface(t *testing.T) {
	t.Parallel()

	require.Implements(t,
		new(validators.OCMError), utils.OCMResponseError(400),
		"utils.OCMResponseError must implement the validators.OCMError interface",
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

			assert.Equal(t, utils.OCMResponseError(tc.code).ServerSide(), tc.result)
		})
	}
}

func TestEnvOCMTokenProviderInterface(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(utils.OCMTokenProvider), utils.EnvOCMTokenProvider{})
}

func TestEnvOCMTokenProviderProvideToken(t *testing.T) {
	t.Parallel()

	const (
		goodEnvVar = "TEST_VAR"
		goodToken  = "supersecrettoken"
	)

	err := os.Setenv(goodEnvVar, goodToken)
	require.NoError(t, err)

	for name, tc := range map[string]struct {
		envVar        string
		expectedToken string
		assertFunc    func(assert.TestingT, error, ...interface{}) bool
	}{
		"correctly set environment variable": {
			envVar:        goodEnvVar,
			expectedToken: goodToken,
			assertFunc:    assert.NoError,
		},
		"unset environment variable": {
			envVar:        "DNE",
			expectedToken: "",
			assertFunc:    assert.Error,
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			token, err := utils.NewEnvOCMTokenProvider(tc.envVar).ProvideToken()
			tc.assertFunc(t, err)
			require.Equal(t, tc.expectedToken, token)
		})
	}
}
