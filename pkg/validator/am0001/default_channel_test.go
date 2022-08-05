package am0001

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"github.com/stretchr/testify/require"
)

func TestDefaultChannelValid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewDefaultChannel)

	validBundles, err := testutils.DefaultValidBundleMap()
	require.NoError(t, err)

	tester.TestValidBundles(validBundles)

}

func TestDefaultChannelInvalid(t *testing.T) {
	t.Parallel()

	tester := testutils.NewValidatorTester(t, NewDefaultChannel)

	tester.TestInvalidBundles(map[string]types.MetaBundle{
		"multiple channels": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator",
				DefaultChannel: "alpha",
				Channels: &[]v1alpha1.Channel{
					{
						Name: "beta",
					},
					{
						Name: "sigma",
					},
				},
			},
		},
		"mismatched channels": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator",
				DefaultChannel: "beta",
				Channels: &[]v1alpha1.Channel{
					{
						Name: "alpha",
					},
				},
			},
		},
		"unknown channel": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator",
				DefaultChannel: "invalid",
			},
		},
		"not present in bundle channels annotation": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator",
				DefaultChannel: "alpha",
				Channels: &[]v1alpha1.Channel{
					{
						Name: "alpha",
					},
				},
			},
			Bundles: []*registry.Bundle{
				{
					Annotations: &registry.Annotations{
						DefaultChannelName: "alpha",
						Channels:           "beta,stable,rc",
					},
				},
			},
		},
		"not present in the bundle default channel": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator",
				DefaultChannel: "alpha",
				Channels: &[]v1alpha1.Channel{
					{
						Name: "alpha",
					},
				},
			},
			Bundles: []*registry.Bundle{
				{
					Annotations: &registry.Annotations{
						DefaultChannelName: "beta",
						Channels:           "beta,stable,rc,alpha",
					},
				},
			},
		},
		"DefaultChannel is not alpha for bundle default channel not defined": {
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				ID:             "random-operator",
				DefaultChannel: "stable",
				Channels: &[]v1alpha1.Channel{
					{
						Name: "stable",
					},
				},
			},
			Bundles: []*registry.Bundle{
				{
					Annotations: &registry.Annotations{
						DefaultChannelName: "",
						Channels:           "stable,beta",
					},
				},
			},
		},
	})
}
