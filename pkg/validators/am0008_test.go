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
	return validators.AM0008.Runner(mb)
}

func (val TestAM0008) SucceedingCandidates() []types.MetaBundle {
	res := testutils.DefaultSucceedingCandidates()
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
	return append(res, moreSucceedingCandidates...)
}

func (val TestAM0008) FailingCandidates() []types.MetaBundle {
	return []types.MetaBundle{
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
}
