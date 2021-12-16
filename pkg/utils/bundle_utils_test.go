package utils_test

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"

	"github.com/stretchr/testify/require"
)

func TestExtractAndParseAddons(t *testing.T) {
	defer testutils.RemoveDir(utils.DefaultDownloadPath)
	testCases := []struct {
		indexImage             string
		operatorName           string
		expectedErrorSubstring *string
	}{
		{
			indexImage:             "quay.io/osd-addons/reference-addon-index@sha256:b9e87a598e7fd6afb4bfedb31e4098435c2105cc8ebe33231c341e515ba9054d",
			operatorName:           "reference-addon",
			expectedErrorSubstring: nil,
		},
		{
			indexImage:             "quay.io/osd-addons/reference-addon-index@sha256:b9e87a598e7fd6afb4bfedb31e4098435c2105cc8ebe33231c341e515ba9054d",
			operatorName:           "lorem-ipsum",
			expectedErrorSubstring: testutils.GetStringLiteralRef("can't find any bundles for the operator 'lorem-ipsum'"),
		},
	}

	for _, tc := range testCases {
		tc := tc // pin
		t.Run(tc.operatorName, func(t *testing.T) {
			t.Parallel()
			bundles, err := utils.ExtractAndParseAddons(tc.indexImage, tc.operatorName)
			if tc.expectedErrorSubstring == nil {
				require.Greater(t, len(bundles), 0)
				require.NoError(t, err)
			} else {
				require.Equal(t, len(bundles), 0)
				require.Contains(t, err.Error(), *tc.expectedErrorSubstring)
			}
		})
	}
}
