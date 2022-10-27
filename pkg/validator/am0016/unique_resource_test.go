package am0016

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	mtsrev1 "github.com/mt-sre/addon-metadata-operator/pkg/mtsre/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
)

func TestUniqueResourceValid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewUniqueResource)
	tester.TestValidBundles(map[string]types.MetaBundle{
		"no additional catalog sources": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AdditionalCatalogSources: nil,
			},
		},
		"one additional catalog source": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AdditionalCatalogSources: &[]mtsrev1.AdditionalCatalogSource{
					{
						Name:  "test-name-1",
						Image: "test-image-1",
					},
				},
			},
		},
		"multiple additional catalog sources": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AdditionalCatalogSources: &[]mtsrev1.AdditionalCatalogSource{
					{
						Name:  "test-name-1",
						Image: "test-image-1",
					},
					{
						Name:  "test-name-2",
						Image: "test-image-2",
					},
					{
						Name:  "test-name-3",
						Image: "test-image-3",
					},
				},
			},
		},
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
		"no credential request": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				CredentialsRequests: nil,
			},
		},
		"one credential request": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				CredentialsRequests: &[]mtsrev1.CredentialsRequest{
					{
						Name: "test-cr",
					},
				},
			},
		},
		"multiple credential requests": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				CredentialsRequests: &[]mtsrev1.CredentialsRequest{
					{
						Name: "test-cr-1",
					},
					{
						Name: "test-cr-2",
					},
					{
						Name: "test-cr-3",
					},
				},
			},
		},
		"additional catalog source and secrets": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AdditionalCatalogSources: &[]mtsrev1.AdditionalCatalogSource{
					{
						Name:  "test-name-1",
						Image: "test-image-1",
					},
					{
						Name:  "test-name-2",
						Image: "test-image-2",
					},
				},
				Config: &mtsrev1.Config{
					Secrets: &[]mtsrev1.Secret{
						{
							Name: "test-secret-1",
						},
						{
							Name: "test-secret-2",
						},
					},
				},
			},
		},
		"additional catalog source and credential requests": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AdditionalCatalogSources: &[]mtsrev1.AdditionalCatalogSource{
					{
						Name:  "test-name-1",
						Image: "test-image-1",
					},
					{
						Name:  "test-name-2",
						Image: "test-image-2",
					},
				},
				CredentialsRequests: &[]mtsrev1.CredentialsRequest{
					{
						Name: "test-cr-1",
					},
					{
						Name: "test-cr-2",
					},
				},
			},
		},
		"secrets and credential requests": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				Config: &mtsrev1.Config{
					Secrets: &[]mtsrev1.Secret{
						{
							Name: "test-secret-1",
						},
						{
							Name: "test-secret-2",
						},
					},
				},
				CredentialsRequests: &[]mtsrev1.CredentialsRequest{
					{
						Name: "test-cr-1",
					},
					{
						Name: "test-cr-2",
					},
				},
			},
		},
		"all three resources": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AdditionalCatalogSources: &[]mtsrev1.AdditionalCatalogSource{
					{
						Name:  "test-name-1",
						Image: "test-image-1",
					},
					{
						Name:  "test-name-2",
						Image: "test-image-2",
					},
				},
				Config: &mtsrev1.Config{
					Secrets: &[]mtsrev1.Secret{
						{
							Name: "test-secret-1",
						},
						{
							Name: "test-secret-2",
						},
					},
				},
				CredentialsRequests: &[]mtsrev1.CredentialsRequest{
					{
						Name: "test-cr-1",
					},
					{
						Name: "test-cr-2",
					},
				},
			},
		},
	})
}

func TestUniqueResourcesInvalid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewUniqueResource)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"one additional catalog source repeating": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AdditionalCatalogSources: &[]mtsrev1.AdditionalCatalogSource{
					{
						Name:  "test-name-1",
						Image: "test-image-1",
					},
					{
						Name:  "test-name-1",
						Image: "test-image-2",
					},
					{
						Name:  "test-name-3",
						Image: "test-image-3",
					},
				},
			},
		},
		"multiple additional catalog source repeating": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AdditionalCatalogSources: &[]mtsrev1.AdditionalCatalogSource{
					{
						Name:  "test-name-1",
						Image: "test-image-1",
					},
					{
						Name:  "test-name-2",
						Image: "test-image-2",
					},
					{
						Name:  "test-name-3",
						Image: "test-image-3",
					},
					{
						Name:  "test-name-2",
						Image: "test-image-4",
					},
					{
						Name:  "test-name-1",
						Image: "test-image-5",
					},
				},
			},
		},
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
		"one credential request repeating": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				CredentialsRequests: &[]mtsrev1.CredentialsRequest{
					{
						Name: "test-cr-1",
					},
					{
						Name: "test-cr-2",
					},
					{
						Name: "test-cr-1",
					},
				},
			},
		},
		"multiple credential requests repeating": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				CredentialsRequests: &[]mtsrev1.CredentialsRequest{
					{
						Name: "test-cr-1",
					},
					{
						Name: "test-cr-2",
					},
					{
						Name: "test-cr-1",
					},
					{
						Name: "test-cr-3",
					},
					{
						Name: "test-cr-2",
					},
				},
			},
		},
		"additional catalog source and secrets repeating": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AdditionalCatalogSources: &[]mtsrev1.AdditionalCatalogSource{
					{
						Name:  "test-name-1",
						Image: "test-image-1",
					},
					{
						Name:  "test-name-1",
						Image: "test-image-2",
					},
					{
						Name:  "test-name-3",
						Image: "test-image-3",
					},
				},
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
		"additional catalog source and credentail requests repeating": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AdditionalCatalogSources: &[]mtsrev1.AdditionalCatalogSource{
					{
						Name:  "test-name-1",
						Image: "test-image-1",
					},
					{
						Name:  "test-name-1",
						Image: "test-image-2",
					},
					{
						Name:  "test-name-3",
						Image: "test-image-3",
					},
				},
				CredentialsRequests: &[]mtsrev1.CredentialsRequest{
					{
						Name: "test-cr-1",
					},
					{
						Name: "test-cr-2",
					},
					{
						Name: "test-cr-1",
					},
				},
			},
		},
		"secrets and credential requests repeating": {
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
				CredentialsRequests: &[]mtsrev1.CredentialsRequest{
					{
						Name: "test-cr-1",
					},
					{
						Name: "test-cr-2",
					},
					{
						Name: "test-cr-1",
					},
				},
			},
		},
		"all resources repeating": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				AdditionalCatalogSources: &[]mtsrev1.AdditionalCatalogSource{
					{
						Name:  "test-name-1",
						Image: "test-image-1",
					},
					{
						Name:  "test-name-1",
						Image: "test-image-2",
					},
					{
						Name:  "test-name-3",
						Image: "test-image-3",
					},
				},
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
						{
							Name: "test-secret-2",
						},
					},
				},
				CredentialsRequests: &[]mtsrev1.CredentialsRequest{
					{
						Name: "test-cr-1",
					},
					{
						Name: "test-cr-2",
					},
					{
						Name: "test-cr-1",
					},
					{
						Name: "test-cr-3",
					},
					{
						Name: "test-cr-2",
					},
				},
			},
		},
	})
}
