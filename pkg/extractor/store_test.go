package extractor

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThreadSafeStoreInterfaces(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(Store), new(ThreadSafeStore))
}

func TestThreadSafeStore(t *testing.T) {
	t.Parallel()

	const numWorkers = 10

	store := NewThreadSafeStore()

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		i := i

		wg.Add(1)

		go func() {
			defer wg.Done()

			err := store.Write(i, i)
			require.NoError(t, err)

			val, ok := store.Read(i)
			require.True(t, ok)

			assert.Equal(t, i, val)
		}()
	}

	wg.Wait()
}
