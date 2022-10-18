package am0017

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	mtsrev1 "github.com/mt-sre/addon-metadata-operator/pkg/mtsre/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
)

func TestSecretValid(t *testing.T) {
	t.Parallel()
	tester := testutils.NewValidatorTester(t, NewSecret)
	tester.TestValidBundles(map[string]types.MetaBundle{
		"config is nil": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				Config: nil,
			},
		},
		"no secret": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				Config: &mtsrev1.Config{
					Secrets: nil,
				},
			},
		},
		"one secret": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				Config: &mtsrev1.Config{
					Secrets: &[]mtsrev1.Secret{
						{
							Name: "test-secret",
						},
					},
				},
			},
		},
		"multiple secrets": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				Config: &mtsrev1.Config{
					Secrets: &[]mtsrev1.Secret{
						{
							Name: "test-secret-1",
						},
						{
							Name: "test-secret-2",
						},
						{
							Name: "test-secret-3",
						},
					},
				},
			},
		},
	})
}

func TestSecretInvalid(t *testing.T) {
	t.Parallel()
	tester := testutils.NewValidatorTester(t, NewSecret)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"one secret repeating": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				Config: &mtsrev1.Config{
					Secrets: &[]mtsrev1.Secret{
						{
							Name: "test-secret-1",
						},
						{
							Name: "test-secret-2",
						},
						{
							Name: "test-secret-1",
						},
					},
				},
			},
		},
		"multiple secret repeating": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				Config: &mtsrev1.Config{
					Secrets: &[]mtsrev1.Secret{
						{
							Name: "test-secret-1",
						},
						{
							Name: "test-secret-2",
						},
						{
							Name: "test-secret-2",
						},
						{
							Name: "test-secret-3",
						},
						{
							Name: "test-secret-3",
						},
					},
				},
			},
		},
	})
}
