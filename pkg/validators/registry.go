package validators

import (
	"log"
	"sort"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

// Registry - holds all registered Validators
var Registry = NewValidatorsRegistry()

func NewValidatorsRegistry() *defaultRegistry {
	return &defaultRegistry{Data: make(map[string]types.Validator)}
}

type defaultRegistry struct {
	Data map[string]types.Validator
}

func (r *defaultRegistry) Add(v types.Validator) {
	if _, ok := r.Data[v.Code]; ok {
		log.Panicf("Validator code %v already exist.", v.Code)
	}
	r.Data[v.Code] = v
}

func (r *defaultRegistry) Len() int {
	return len(r.Data)
}

func (r *defaultRegistry) All() map[string]types.Validator {
	return r.Data
}

func (r *defaultRegistry) Get(k string) (types.Validator, bool) {
	v, ok := r.Data[k]
	return v, ok
}

func (r *defaultRegistry) ListSorted() types.ValidatorList {
	res := make(types.ValidatorList, 0)
	for _, v := range r.Data {
		res = append(res, v)
	}
	sort.Sort(res)
	return res
}
