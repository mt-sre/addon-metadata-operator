package v1

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestAddOnRequirementDataDynamicTyping(t *testing.T) {
	b := []byte(`{
		"id": "dummy",
		"name": "dummy",
		"description": "dummy",
		"value_type": "string",
		"required": true,
		"editable": true,
		"enabled": true,
		"conditions": [
			{
				"resource": "addon",
				"data": {
					"product.id": "aws",
					"cloud_provider.id": ["aws"],
				},
			},
		],
	}`)

	var param AddOnParameter
	err := yaml.Unmarshal(b, &param)
	require.NoError(t, err)

	for _, req := range *param.Conditions {
		for _, key := range []string{"product.id", "cloud_provider.id"} {
			valString, valSlice, err := stringOrSlice(req.Data[key].Raw)
			require.NoError(t, err)
			switch key {
			case "product.id":
				require.Nil(t, valSlice)
				require.Equal(t, *valString, "aws")
			case "cloud_provider.id":
				require.Nil(t, valString)
				require.Equal(t, *valSlice, []string{"aws"})
			}
		}
	}
}

func stringOrSlice(b []byte) (*string, *[]string, error) {
	var resString string
	var resSlice []string
	if err := json.Unmarshal(b, &resString); err == nil {
		return &resString, nil, nil
	}
	if err := json.Unmarshal(b, &resSlice); err == nil {
		return nil, &resSlice, nil
	}
	return nil, nil, errors.New("invalid type")
}
