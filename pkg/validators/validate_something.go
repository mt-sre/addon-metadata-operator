package validators

import (
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

func ValidateSomethingCrossBetweenImageSetAndMetadata(metabundle *utils.MetaBundle) (bool, error) {
	if metabundle.AddonMeta.IndexImage != "" {
		return false, nil
	}
	return true, nil
}
