package validators_test

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
)

func init() {
	TestRegistry.Add(TestAM0002{})
}

type TestAM0002 struct{}

func (val TestAM0002) Name() string {
	return validators.AM0002.Name
}

func (val TestAM0002) Run(mb types.MetaBundle) types.ValidatorResult {
	return validators.AM0002.Runner(mb)
}

func (val TestAM0002) SucceedingCandidates() ([]types.MetaBundle, error) {
	res := []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:    "random-operator",
				Label: "api.openshift.com/addon-random-operator",
			},
		},
	}
	return res, nil
}

func (val TestAM0002) FailingCandidates() ([]types.MetaBundle, error) {
	res := []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:    "random-operator",
				Label: "foo-bar",
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:    "random-operator",
				Label: "api.openshift.com/addon-random-operator-x",
			},
		},
	}
	return res, nil
}
