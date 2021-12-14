package utils

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/internal/testutils"

	"github.com/stretchr/testify/require"
)

func TestExtractAndParseAddons(t *testing.T) {
	defer testutils.RemoveDir(DefaultDownloadPath)
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

	for _, testCase := range testCases {
		_, err := ExtractAndParseAddons(testCase.indexImage, testCase.operatorName)
		if testCase.expectedErrorSubstring == nil {
			require.NoError(t, err)
		} else {
			require.Contains(t, err.Error(), *testCase.expectedErrorSubstring)
		}
	}
}
