package validators_test

import (
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
)

func init() {
	TestRegistry.Add(TestAM0012{})
}

type TestAM0012 struct{}

func (val TestAM0012) Name() string {
	return validators.AM0012.Name
}

func (val TestAM0012) Run(mb types.MetaBundle) types.ValidatorResult {
	return validators.AM0012.Runner(mb)
}

func (val TestAM0012) SucceedingCandidates() ([]types.MetaBundle, error) {
	res, err := testutils.DefaultSucceedingCandidates()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (val TestAM0012) FailingCandidates() ([]types.MetaBundle, error) {
	invalidBundles, err := utils.ExtractAndParseAddons(
		"quay.io/osd-addons/rhods-index@sha256:487e106059aea611af377985e6f30d7879bc36c4a16fe0f70531b7c1befd4675",
		"rhods-operator",
	)
	if err != nil {
		return nil, err
	}

	res := []types.MetaBundle{
		{
			Bundles: invalidBundles,
		},
	}
	return res, nil
}
