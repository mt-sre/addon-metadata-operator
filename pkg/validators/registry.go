package validators

import (
	"log"
	"sync"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

// Registry - holds all registered Validators
var Registry = NewDefaultRegistry()

func NewDefaultRegistry() registry {
	return registry{
		&sync.Mutex{},
		make(map[string]utils.Validator),
	}
}

type registry struct {
	*sync.Mutex
	Data map[string]utils.Validator
}

// Add - update the registry in a thread-safe way. Called in init() functions
func (r *registry) Add(v utils.Validator) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.Data[v.Code]; ok {
		log.Fatalf("Validator code %v already exist.", v.Code)
	}
	r.Data[v.Code] = v
}

func (r registry) Len() int {
	return len(r.Data)
}

func (r registry) All() map[string]utils.Validator {
	return r.Data
}

func (r registry) Get(k string) (utils.Validator, bool) {
	v, ok := r.Data[k]
	return v, ok
}

// TestRegistry - register all test structs
var TestRegistry = NewTestRegistry()

func NewTestRegistry() *testRegistry {
	return &testRegistry{
		&sync.Mutex{},
		[]utils.ValidatorTest{},
	}
}

type testRegistry struct {
	*sync.Mutex
	Data []utils.ValidatorTest
}

// Add - update the test registry in a thread-safe way. Called in init() functions
func (t *testRegistry) Add(v utils.ValidatorTest) {
	t.Lock()
	defer t.Unlock()

	t.Data = append(t.Data, v)
}

func (t *testRegistry) All() []utils.ValidatorTest {
	return t.Data
}
