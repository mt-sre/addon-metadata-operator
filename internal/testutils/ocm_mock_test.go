package testutils

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
	"github.com/stretchr/testify/require"
)

func TestMockOCMClientInterface(t *testing.T) {
	t.Parallel()

	require.Implements(t,
		new(validators.OCMClient), new(MockOCMClient),
		"MockOCMClient must implement the types.OCMClient interface",
	)
}
