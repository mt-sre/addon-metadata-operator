package validators

import (
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

func ValidateDefaultChannel(metabundle *utils.MetaBundle) (bool, error) {
	defaultChannel := metabundle.AddonMeta.DefaultChannel
	for _, channel := range metabundle.AddonMeta.Channels {
		if channel.Name == defaultChannel {
			return true, nil
		}
	}
	return false, nil
}
