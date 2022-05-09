package testutils

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	"github.com/stretchr/testify/require"
)

func TestMockOCMClientInterface(t *testing.T) {
	t.Parallel()

	require.Implements(t,
		new(validator.OCMClient), new(MockOCMClient),
		"MockOCMClient must implement the types.OCMClient interface",
	)
}
