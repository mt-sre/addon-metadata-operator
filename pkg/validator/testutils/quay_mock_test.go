package testutils

import (
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"

	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockQuayClientInterfaces(t *testing.T) {
	require.Implements(t, new(validator.QuayClient), new(MockQuayClient))
}
