package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetSemver(t *testing.T) {
	cases := []struct {
		name           string
		expectedSemver string
		isError        bool
	}{
		{
			name:           "reference-addon.v0.0.1",
			expectedSemver: "0.0.1",
			isError:        false,
		},
		{
			name:           "reference-addon.v2.3.2",
			expectedSemver: "2.3.2",
			isError:        false,
		},
		{
			name:           "invalid-semver.v2.3.2.4.5",
			expectedSemver: "",
			isError:        true,
		},
		{
			name:           "invalid_name",
			expectedSemver: "",
			isError:        true,
		},
	}
	for _, tc := range cases {
		tc := tc // pin
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			imageSet := &AddonImageSetSpec{Name: tc.name}
			semver, err := imageSet.GetSemver()
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, semver, tc.expectedSemver)
			}
		})
	}
}
