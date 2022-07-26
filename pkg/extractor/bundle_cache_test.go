package extractor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBundleCacheImplInterfaces(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(BundleCache), new(BundleCacheImpl))
}
