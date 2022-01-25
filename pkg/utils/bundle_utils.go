package utils

import (
	"fmt"

	"github.com/operator-framework/operator-registry/pkg/registry"
)

// GetBundleNameVersion - useful for validation error reporting
func GetBundleNameVersion(b *registry.Bundle) (string, error) {
	version, err := b.Version()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v:%v", b.Name, version), nil
}
