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
	TestRegistry.Add(TestAM0003{})
}

type TestAM0003 struct{}

func (v TestAM0003) Name() string {
	return validators.AM0003.Name
}

func (v TestAM0003) Run(mb types.MetaBundle) types.ValidatorResult {
	return validators.AM0003.Runner(mb)
}

func (v TestAM0003) SucceedingCandidates() ([]types.MetaBundle, error) {
	return testutils.DefaultSucceedingCandidates()

}

func (v TestAM0003) FailingCandidates() ([]types.MetaBundle, error) {
	var res []types.MetaBundle

	type failingMetaBundleGen func() (types.MetaBundle, error)

	for _, fn := range []failingMetaBundleGen{
		metaBundleWithBadPackageNameAnnotation,
		metaBundleWithBadSemver,
		metaBundleWithBadCSVName,
		metaBundleWithBadCSVReplaces,
	} {
		mb, err := fn()
		if err != nil {
			return nil, err
		}
		res = append(res, mb)
	}
	return res, nil
}

func metaBundleWithBadPackageNameAnnotation() (types.MetaBundle, error) {
	return newAM0003FailingMetaBundle("csv_valid.yml", "invalid")
}

func metaBundleWithBadSemver() (types.MetaBundle, error) {
	return newAM0003FailingMetaBundle("csv_semver_invalid.yml", "reference-addon")
}

func metaBundleWithBadCSVName() (types.MetaBundle, error) {
	return newAM0003FailingMetaBundle("csv_name_invalid.yml", "reference-addon")
}

func metaBundleWithBadCSVReplaces() (types.MetaBundle, error) {
	return newAM0003FailingMetaBundle("csv_replaces_invalid.yml", "reference-addon")
}

func newAM0003FailingMetaBundle(csvFile string, pkgNameAnnotation string) (types.MetaBundle, error) {
	bundle, err := newAM0003FailingBundle(csvFile, pkgNameAnnotation)
	if err != nil {
		return types.MetaBundle{}, err
	}

	mb := types.MetaBundle{
		AddonMeta: &v1alpha1.AddonMetadataSpec{
			OperatorName: "reference-addon",
		},
		Bundles: []*registry.Bundle{
			bundle,
		},
	}
	return mb, nil
}

func newAM0003FailingBundle(csvFile, packageName string) (*registry.Bundle, error) {
	csvPath := path.Join(testutils.TestdataDir(), "assets/am0003/", csvFile)
	res, err := testutils.NewBundle("am0003-failing-bundle", csvPath)
	if err != nil {
		return nil, err
	}
	res.Annotations.PackageName = packageName
	return res, nil
}
