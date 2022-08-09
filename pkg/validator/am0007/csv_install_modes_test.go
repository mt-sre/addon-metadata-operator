package am0007

import (
	"path/filepath"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/stretchr/testify/require"
)

func TestCSVInstallModeValid(t *testing.T) {
	t.Parallel()

	loader := testutils.NewBundlerLoader(t)

	bundles, err := testutils.DefaultValidBundleMap()
	require.NoError(t, err)

	for name, bundle := range map[string]types.MetaBundle{
		"OwnNamespace": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "OwnNamespace",
			},
			Bundles: []*registry.Bundle{
				loader.LoadFromCSV(filepath.Join("test_csvs", "csv.yaml")),
			},
		},
		"AllNamespaces": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "AllNamespaces",
			},
			Bundles: []*registry.Bundle{
				loader.LoadFromCSV(filepath.Join("test_csvs", "csv.yaml")),
			},
		},
	} {
		bundles[name] = bundle
	}

	tester := testutils.NewValidatorTester(t, NewCSVInstallModes)
	tester.TestValidBundles(bundles)
}

func TestCSVInstallModeInvalid(t *testing.T) {
	t.Parallel()

	loader := testutils.NewBundlerLoader(t)

	tester := testutils.NewValidatorTester(t, NewCSVInstallModes)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"invalid install mode": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "something",
			},
			Bundles: []*registry.Bundle{
				loader.LoadFromCSV(filepath.Join("test_csvs", "csv.yaml")),
			},
		},
		"invalid install mode/SingleNamespace": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "SingleNamespace",
			},
			Bundles: []*registry.Bundle{
				loader.LoadFromCSV(filepath.Join("test_csvs", "csv.yaml")),
			},
		},
		"invalid install mode/MultiNamespace": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				InstallMode: "MultiNamespace",
			},
			Bundles: []*registry.Bundle{
				loader.LoadFromCSV(filepath.Join("test_csvs", "csv.yaml")),
			},
		},
	})
}
