package am0003

import (
	"path/filepath"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/operator"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
)

func TestOperatorNameValid(t *testing.T) {
	t.Parallel()

	loader := testutils.NewBundlerLoader(t)

	tester := testutils.NewValidatorTester(t, NewOperatorName)
	tester.TestValidBundles(map[string]types.MetaBundle{
		"valid package name annotation": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []operator.Bundle{
				loader.LoadFromCSV(
					filepath.Join("test_csvs", "csv_valid.yaml"),
					testutils.WithPackageName("reference-addon"),
				),
			},
		},
	})
}

func TestOperatornameInvalid(t *testing.T) {
	t.Parallel()

	loader := testutils.NewBundlerLoader(t)

	tester := testutils.NewValidatorTester(t, NewOperatorName)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"invalid package name annotation": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []operator.Bundle{
				loader.LoadFromCSV(
					filepath.Join("test_csvs", "csv_valid.yaml"),
					testutils.WithPackageName("invalid"),
				),
			},
		},
		"invalid csv name identifier": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []operator.Bundle{
				loader.LoadFromCSV(
					filepath.Join("test_csvs", "csv_name_invalid.yaml"),
					testutils.WithPackageName("reference-addon"),
				),
			},
		},
		"invalid replaces identifier": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []operator.Bundle{
				loader.LoadFromCSV(
					filepath.Join("test_csvs", "csv_replaces_invalid.yaml"),
					testutils.WithPackageName("reference-addon"),
				),
			},
		},
		"invalid semver": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []operator.Bundle{
				loader.LoadFromCSV(
					filepath.Join("test_csvs", "csv_semver_invalid.yaml"),
					testutils.WithPackageName("reference-addon"),
				),
			},
		},
	})
}
