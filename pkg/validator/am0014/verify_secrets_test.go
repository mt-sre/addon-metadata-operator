package am0014

import (
	_ "embed"
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	mtsrev1 "github.com/mt-sre/addon-metadata-operator/pkg/mtsre/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/stretchr/testify/require"
)

func TestVerifySecretParamsValid(t *testing.T) {
	t.Parallel()

	bundles, err := testutils.DefaultValidBundleMap()
	require.NoError(t, err)

	for name, bundle := range map[string]types.MetaBundle{
		"valid secrets defined": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "random-operator-1",
				Secrets: &[]mtsrev1.Secret{
					{
						Name:      "secret-one",
						Type:      "kubernetes.io/dockerconfigjson",
						VaultPath: "mtsre/quay/osd-addons/secrets/random-operator-1/secret-one",
					},
					{
						Name:      "secret-two",
						Type:      "bootstrap.kubernetes.io/token",
						VaultPath: "mtsre/quay/osd-addons/secrets/random-operator-1/secret-two",
					},
				},
			},
		},
		"no secrets defined": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "random-operator-2",
			},
		},
	} {
		bundles[name] = bundle
	}

	tester := testutils.NewValidatorTester(t, NewVerifySecretParams)
	tester.TestValidBundles(bundles)
}

func TestVerifySecretParamsInvalid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewVerifySecretParams)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"secrets not present in deploy.yaml": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID: "random-operator-1",
				Secrets: &[]mtsrev1.Secret{
					{
						Name:      "secretsss-1",
						Type:      "kubernetes.io/dockerconfigjson",
						VaultPath: "path/to/secret",
					},
					{
						Name:      "secret99",
						Type:      "bootstrap.kubernetes.io/token",
						VaultPath: "test/path/to/some/secret",
					},
					{
						Name:      "READER_ENDPOINT",
						Type:      "kubernetes.io/dockerconfigjson",
						VaultPath: "random/wrong/path/here",
					},
				},
			},
		},
	})
}
