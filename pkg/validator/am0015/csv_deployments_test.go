package am0015

import (
	"path/filepath"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/operator"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
)

func TestCSVDeploymentValid(t *testing.T) {
	t.Parallel()

	loader := testutils.NewBundlerLoader(t)

	tester := testutils.NewValidatorTester(t, NewCSVDeployment)
	tester.TestValidBundles(map[string]types.MetaBundle{
		"valid CSV Deployment": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []operator.Bundle{
				loader.LoadFromCSV(
					filepath.Join("test_csvs", "csv_valid.yaml"),
				),
			},
		},
	})
}

func TestCSVDeploymentInvalid(t *testing.T) {
	t.Parallel()

	loader := testutils.NewBundlerLoader(t)

	tester := testutils.NewValidatorTester(t, NewCSVDeployment)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"invalid csv livenessprobe": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []operator.Bundle{
				loader.LoadFromCSV(
					filepath.Join("test_csvs", "csv_invalid_livenessprobe.yaml"),
				),
			},
		},
		"invalid csv readinessprobe": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []operator.Bundle{
				loader.LoadFromCSV(
					filepath.Join("test_csvs", "csv_invalid_readinessprobe.yaml"),
				),
			},
		},
		"invalid csv memory requirement": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []operator.Bundle{
				loader.LoadFromCSV(
					filepath.Join("test_csvs", "csv_invalid_memory_requirement.yaml"),
				),
			},
		},
		"invalid csv cpu requirements": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []operator.Bundle{
				loader.LoadFromCSV(
					filepath.Join("test_csvs", "csv_invalid_cpu_requirement.yaml"),
				),
			},
		},
	})
}
