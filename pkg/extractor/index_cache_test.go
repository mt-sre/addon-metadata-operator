package extractor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndexCacheImplInterfaces(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(IndexCache), new(IndexCacheImpl))
}
