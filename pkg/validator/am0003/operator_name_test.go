package am0003

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

func TestOperatorNameValid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewOperatorName)
	tester.TestValidBundles(map[string]types.MetaBundle{
		"valid package name annotation": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []*registry.Bundle{
				newTestBundle(t, "csv_valid.yaml", "reference-addon"),
			},
		},
	})
}

func TestOperatornameInvalid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewOperatorName)

	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"invalid package name annotation": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []*registry.Bundle{
				newTestBundle(t, "csv_valid.yaml", "invalid"),
			},
		},
		"invalid csv name identifier": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []*registry.Bundle{
				newTestBundle(t, "csv_name_invalid.yaml", "reference-addon"),
			},
		},
		"invalid replaces identifier": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []*registry.Bundle{
				newTestBundle(t, "csv_replaces_invalid.yaml", "reference-addon"),
			},
		},
		"invalid semver": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				OperatorName: "reference-addon",
			},
			Bundles: []*registry.Bundle{
				newTestBundle(t, "csv_semver_invalid.yaml", "reference-addon"),
			},
		},
	})
}

func newTestBundle(t *testing.T, csvName, packageName string) *registry.Bundle {
	t.Helper()

	res, err := utils.NewBundle("am0003", filepath.Join(utils.RootDir().TestData().Validators(), "am0003", csvName))
	require.NoError(t, err)

	res.Annotations.PackageName = packageName

	return res
}
