package am0008

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	utils "github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/stretchr/testify/require"
)

func TestNamespaceValid(t *testing.T) {
	t.Parallel()

	bundles, err := utils.DefaultValidBundleMap()
	require.NoError(t, err)

	for name, bundle := range map[string]types.MetaBundle{
		"redhat prefixed namespace": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:              "random-operator",
				TargetNamespace: "redhat-random-operator",
				Namespaces: []string{
					"redhat-random-operator",
				},
			},
		},
	} {
		bundles[name] = bundle
	}

	tester := utils.NewValidatorTester(t, NewNamespace)
	tester.TestValidBundles(bundles)

}

func TestNamespaceInvalid(t *testing.T) {
	t.Parallel()

	tester := utils.NewValidatorTester(t, NewNamespace)
	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"targetNamespace not in Namespaces": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:              "random-operator",
				TargetNamespace: "redhat-other-operator",
				Namespaces: []string{
					"redhat-random-operator",
				},
			},
		},
		"non redhat prefixed namespace": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:              "random-operator-1",
				TargetNamespace: "redhat-random-operator",
				Namespaces: []string{
					"redhat-random-operator",
					"other-operator",
				},
			},
		},
	})
}
