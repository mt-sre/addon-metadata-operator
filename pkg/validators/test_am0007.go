package validators

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
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

func (v TestAM0007) SucceedingCandidates() []utils.MetaBundle {
	csvObj, err := yamlToDynamicObj("../validators/assets/am0007/csv.yaml") // from the pov of validate/validate_test.go
	if err != nil {
		panic(err)
	}

	succeedingCandidates := []utils.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "OwnNamespace",
			},
			Bundles: []registry.Bundle{
				*registry.NewBundle("random-bundle", &registry.Annotations{}, &csvObj),
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "AllNamespaces",
			},
			Bundles: []registry.Bundle{
				*registry.NewBundle("random-bundle", &registry.Annotations{}, &csvObj),
			},
		},
	}
	return append(succeedingCandidates, testutils.DefaultSucceedingCandidates()...)
}

// not implemented
func (v TestAM0007) FailingCandidates() []utils.MetaBundle {
	csvObj, err := yamlToDynamicObj("../validators/assets/am0007/csv.yaml")
	if err != nil {
		panic(err)
	}

	return []utils.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "something",
			},
			Bundles: []registry.Bundle{
				*registry.NewBundle("random-bundle", &registry.Annotations{}, &csvObj),
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "SingleNamespace",
			},
			Bundles: []registry.Bundle{
				*registry.NewBundle("random-bundle", &registry.Annotations{}, &csvObj),
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "MultiNamespace",
			},
			Bundles: []registry.Bundle{
				*registry.NewBundle("random-bundle", &registry.Annotations{}, &csvObj),
			},
		},
	}
}
