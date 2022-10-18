package am0016

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	mtsrev1 "github.com/mt-sre/addon-metadata-operator/pkg/mtsre/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
)

func TestAdditionalCatalogSourceValid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewAdditionalCatalogSource)
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
	})
}

func TestAdditionalCatalogSourceInvalid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewAdditionalCatalogSource)
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
	})
}
