package am0007

import (
	"path/filepath"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	utils "github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/stretchr/testify/require"
)

func TestCSVInstallModeValid(t *testing.T) {
	t.Parallel()

	bundles, err := utils.DefaultValidBundleMap()
	require.NoError(t, err)

	for name, bundle := range map[string]types.MetaBundle{
		"OwnNamespace": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "OwnNamespace",
			},
			Bundles: []*registry.Bundle{
				getTestBundle(t),
			},
		},
		"AllNamespaces": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "AllNamespaces",
			},
			Bundles: []*registry.Bundle{
				getTestBundle(t),
			},
		},
	} {
		bundles[name] = bundle
	}

	tester := utils.NewValidatorTester(t, NewCSVInstallModes)
	tester.TestValidBundles(bundles)
}

func TestCSVInstallModeInvalid(t *testing.T) {
	t.Parallel()

	tester := utils.NewValidatorTester(t, NewCSVInstallModes)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"invalid install mode": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "something",
			},
			Bundles: []*registry.Bundle{
				getTestBundle(t),
			},
		},
		"invalid install mode/SingleNamespace": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "SingleNamespace",
			},
			Bundles: []*registry.Bundle{
				getTestBundle(t),
			},
		},
		"invalid install mode/MultiNamespace": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "MultiNamespace",
			},
			Bundles: []*registry.Bundle{
				getTestBundle(t),
			},
		},
	})
}

func getTestBundle(t *testing.T) *registry.Bundle {
	t.Helper()

	res, err := testutils.NewBundle("random-bundle", filepath.Join(testutils.RootDir().TestData().Validators(), "am0007", "csv.yaml"))
	require.NoError(t, err)

	return res
}
