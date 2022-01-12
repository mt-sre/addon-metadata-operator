package validators

import (
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

func init() {
	TestRegistry.Add(TestAM0007{})
}

type TestAM0007 struct{}

func (v TestAM0007) Name() string {
	return AM0007.Name
}

func (v TestAM0007) Run(mb types.MetaBundle) types.ValidatorResult {
	return AM0007.Runner(mb)
}

func (v TestAM0007) SucceedingCandidates() []types.MetaBundle {
	return testutils.DefaultSucceedingCandidates()
}

// not implemented
func (v TestAM0007) FailingCandidates() []types.MetaBundle {
	return []types.MetaBundle{}
}
