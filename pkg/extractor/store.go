package extractor

import "sync"

// Store provides object storage for arbitrary keys.
type Store interface {
	// Read retrieves data with a given id. 'ok'
	// returns false if the id does not exist.
	Read(id interface{}) (data interface{}, ok bool)
	// Write sets the given data for the given id.
	// An error should be returned if the data
	// could not be applied successfully.
	Write(id, data interface{}) error
}

// NewThreadSafeStore returns an implementation of the
// 'Store' interface which is safe for concurrent use
// by multiple goroutines.
func NewThreadSafeStore() *ThreadSafeStore {
	return &ThreadSafeStore{
		store: make(map[interface{}]interface{}),
	}
}

type ThreadSafeStore struct {
	lock  sync.RWMutex
	store map[interface{}]interface{}
}

func (s *ThreadSafeStore) Read(id interface{}) (interface{}, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	data, ok := s.store[id]

	return data, ok
}

func (s *ThreadSafeStore) Write(id, data interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.store[id] = data

	return nil
}
