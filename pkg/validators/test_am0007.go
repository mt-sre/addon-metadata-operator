package validators

import (
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

func init() {
	TestRegistry.Add(TestAM0007{})
}

type TestAM0007 struct{}

func (v TestAM0007) Name() string {
	return AM0007.Name
}

func (v TestAM0007) Run(mb utils.MetaBundle) (bool, string, error) {
	return AM0007.Runner(mb)
}

func (v TestAM0007) SucceedingCandidates() []utils.MetaBundle {
	return testutils.DefaultSucceedingCandidates()
}

// not implemented
func (v TestAM0007) FailingCandidates() []utils.MetaBundle {
	return []utils.MetaBundle{}
}
