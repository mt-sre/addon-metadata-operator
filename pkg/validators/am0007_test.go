package validators_test

import (
	"path"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

func init() {
	TestRegistry.Add(TestAM0007{})
}

type TestAM0007 struct{}

func (v TestAM0007) Name() string {
	return validators.AM0007.Name
}

func (v TestAM0007) Run(mb types.MetaBundle) types.ValidatorResult {
	return validators.AM0007.Validate(mb)
}

func (v TestAM0007) SucceedingCandidates() ([]types.MetaBundle, error) {
	testBundle, err := loadAM0007TestBundle()
	if err != nil {
		return nil, err
	}

	res, err := testutils.DefaultSucceedingCandidates()
	if err != nil {
		return nil, err
	}

	moreSucceedingCandidates := []types.MetaBundle{
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
	return append(moreSucceedingCandidates, res...), nil
}

func (v TestAM0007) FailingCandidates() ([]types.MetaBundle, error) {
	testBundle, err := loadAM0007TestBundle()
	if err != nil {
		return nil, err
	}

	res := []types.MetaBundle{
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
	return res, nil
}

func loadAM0007TestBundle() (registry.Bundle, error) {
	return testutils.NewBundle("random-bundle", path.Join(testutils.TestdataDir(), "assets/am0007/csv.yaml"))
}
