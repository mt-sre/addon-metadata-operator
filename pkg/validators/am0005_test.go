package validators_test

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
)

func init() {
	TestRegistry.Add(TestAM0005{})
}

type TestAM0005 struct{}

func (val TestAM0005) Name() string {
	return validators.AM0005.Name
}

func (val TestAM0005) Run(mb types.MetaBundle) types.ValidatorResult {
	return validators.AM0005.Runner(mb)
}

func (val TestAM0005) SucceedingCandidates() ([]types.MetaBundle, error) {
	res, err := testutils.DefaultSucceedingCandidates()
	if err != nil {
		return nil, err
	}

	moreSucceedingCandidates := []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "quay.io/miwilson/addon-samples",
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "quay.io/asnaraya/reference-addon-test-harness:fix",
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "quay.io/asnaraya/reference-addon-test-harness@sha256:bdc32a600202d36fec4524dbec177e9313ef82ad4bda5bd24d4b75236ca8a482",
			},
		},
	}
	return append(res, moreSucceedingCandidates...), nil
}

func (val TestAM0005) FailingCandidates() ([]types.MetaBundle, error) {
	res := []types.MetaBundle{
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "abcd",
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "quay.io/asnaraya/reference-addon-test-harness:404",
			},
		},
		{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:          "random-operator",
				TestHarness: "https://docker.io/ashishmax31/addon-operator-bundle:0.1.0-cb328d9",
			},
		},
	}
	return res, nil
}
