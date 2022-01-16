package validators_test

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
)

func init() {
	TestRegistry.Add(TestAM0008{})
}

type TestAM0008 struct{}

func (val TestAM0008) Name() string {
	return validators.AM0008.Name
}

func (val TestAM0008) Run(mb types.MetaBundle) types.ValidatorResult {
	return validators.AM0008.Validate(mb)
}

func (val TestAM0008) SucceedingCandidates() ([]types.MetaBundle, error) {
	res, err := testutils.DefaultSucceedingCandidates()
	if err != nil {
		return nil, err
	}
	moreSucceedingCandidates := []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:              "random-operator",
				TargetNamespace: "redhat-random-operator",
				Namespaces: []string{
					"redhat-random-operator",
				},
			},
		},
	}
	return append(res, moreSucceedingCandidates...), nil
}

func (val TestAM0008) FailingCandidates() ([]types.MetaBundle, error) {
	res := []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:              "random-operator",
				TargetNamespace: "redhat-other-operator",
				Namespaces: []string{
					"redhat-random-operator",
				},
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:              "random-operator-1",
				TargetNamespace: "redhat-random-operator",
				Namespaces: []string{
					"redhat-random-operator",
					"other-operator",
				},
			},
		},
	}
	return res, nil
}
