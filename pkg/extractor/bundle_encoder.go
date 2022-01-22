package extractor

import (
	"encoding/json"

	"github.com/golang/snappy"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

type JSONSnappyEncoder struct{}

// NewJSONSnappyEncoder - encodes a bundle to snappy compressed JSON
func NewJSONSnappyEncoder() BundleEncoder {
	return JSONSnappyEncoder{}
}

func (e JSONSnappyEncoder) Encode(bundle *registry.Bundle) ([]byte, error) {
	data, err := json.Marshal(bundle)
	if err != nil {
		return nil, err
	}
	return snappy.Encode(nil, data), nil
}

func (e JSONSnappyEncoder) Decode(data []byte) (*registry.Bundle, error) {
	rawJSON, err := snappy.Decode(nil, data)
	if err != nil {
		return nil, err
	}
	bundle := registry.Bundle{}
	if err := json.Unmarshal(rawJSON, &bundle); err != nil {
		return nil, err
	}
	return &bundle, nil
}
