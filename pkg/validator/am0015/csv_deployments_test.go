package am0015

import (
	"path/filepath"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	utils "github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/stretchr/testify/require"
)

func TestCSVDeploymentValid(t *testing.T) {
	t.Parallel()
	tester := testutils.NewValidatorTester(t, NewCSVDeployment)
	tester.TestValidBundles(map[string]types.MetaBundle{
		"valid CSV Deployment": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []*registry.Bundle{
				newTestBundle(t, "csv_valid.yaml", "reference-addon"),
			},
		},
	})
}

func TestCSVDeploymentInvalid(t *testing.T) {
	t.Parallel()
	tester := testutils.NewValidatorTester(t, NewCSVDeployment)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"invalid csv livenessprobe": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []*registry.Bundle{
				newTestBundle(t, "csv_invalid_livenessprobe.yaml", "reference-addon"),
			},
		},
		"invalid csv readinessprobe": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []*registry.Bundle{
				newTestBundle(t, "csv_invalid_readinessprobe.yaml", "reference-addon"),
			},
		},
		"invalid csv memory requirement": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []*registry.Bundle{
				newTestBundle(t, "csv_invalid_memory_requirement.yaml", "reference-addon"),
			},
		},
		"invalid csv cpu requirements": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []*registry.Bundle{
				newTestBundle(t, "csv_invalid_cpu_requirement.yaml", "reference-addon"),
			},
		},
	})
}

func newTestBundle(t *testing.T, csvName, packageName string) *registry.Bundle {
	t.Helper()

	res, err := utils.NewBundle("am0015", filepath.Join(utils.RootDir().TestData().Validators(), "am0015", csvName))
	require.NoError(t, err)

	res.Annotations.PackageName = packageName

	return res
}
