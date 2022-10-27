package am0017

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	mtsrev1 "github.com/mt-sre/addon-metadata-operator/pkg/mtsre/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
)

func TestNewPullSecretnameValid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewPullSecretname)
	tester.TestValidBundles(map[string]types.MetaBundle{
		"pullSecretName is not present and config is nil": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				PullSecretName: "",
				Config:         nil,
			},
		},
		"pullSecretname is not present and secrets are nil": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				PullSecretName: "",
				Config: &mtsrev1.Config{
					Secrets: nil,
				},
			},
		},
		"pullSecretName is present and addon has one secret": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				PullSecretName: "test-pull-secret",
				Config: &mtsrev1.Config{
					Secrets: &[]mtsrev1.Secret{
						{
							Name: "test-pull-secret",
						},
					},
				},
			},
		},
		"pullSecretName is present and addon has multiple secrets defined": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				PullSecretName: "test-pull-secret-2",
				Config: &mtsrev1.Config{
					Secrets: &[]mtsrev1.Secret{
						{
							Name: "test-pull-secret-1",
						},
						{
							Name: "test-pull-secret-2",
						},
						{
							Name: "test-pull-secret-2",
						},
					},
				},
			},
		},
	})
}

func TestNewPullSecretnameInvalid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewPullSecretname)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"pullSecretName is not nil but addon Config is nil": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				PullSecretName: "test-pull-secret",
				Config:         nil,
			},
		},
		"pullSecretName is not nil but addon secrets are nil": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				PullSecretName: "test-pull-secret",
				Config: &mtsrev1.Config{
					Secrets: nil,
				},
			},
		},
		"pullSecretName is not nil but not present secret": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				PullSecretName: "test-pull-secret",
				Config: &mtsrev1.Config{
					Secrets: &[]mtsrev1.Secret{
						{
							Name: "test-pull",
						},
						{
							Name: "alpha",
						},
						{
							Name: "beta",
						},
					},
				},
			},
		},
	})
}
