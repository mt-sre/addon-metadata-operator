package validators

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/operator-framework/operator-registry/pkg/registry"
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

var (
	testBundle = testutils.NewBundle("random-bundle", "../../internal/testdata/assets/am0007/csv.yaml")
)

func (v TestAM0007) SucceedingCandidates() []types.MetaBundle {
	succeedingCandidates := []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "OwnNamespace",
			},
			Bundles: []registry.Bundle{
				testBundle,
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "AllNamespaces",
			},
			Bundles: []registry.Bundle{
				testBundle,
			},
		},
	}
	return append(succeedingCandidates, testutils.DefaultSucceedingCandidates()...)
}

func (v TestAM0007) FailingCandidates() []types.MetaBundle {
	return []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "something",
			},
			Bundles: []registry.Bundle{
				testBundle,
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "SingleNamespace",
			},
			Bundles: []registry.Bundle{
				testBundle,
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "MultiNamespace",
			},
			Bundles: []registry.Bundle{
				testBundle,
			},
		},
	}
}
