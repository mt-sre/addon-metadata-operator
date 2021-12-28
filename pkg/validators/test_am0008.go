package validators

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

func init() {
	TestRegistry.Add(TestAM0008{})
}

type TestAM0008 struct{}

func (val TestAM0008) Name() string {
	return AM0008.Name
}

func (val TestAM0008) Run(mb utils.MetaBundle) (bool, string, error) {
	return AM0008.Runner(mb)
}

func (val TestAM0008) SucceedingCandidates() []utils.MetaBundle {
	res := testutils.DefaultSucceedingCandidates()
	moreSucceedingCandidates := []utils.MetaBundle{
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

func (val TestAM0008) FailingCandidates() []utils.MetaBundle {
	return []utils.MetaBundle{
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
