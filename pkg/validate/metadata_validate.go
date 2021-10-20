package validate

import "fmt"

func ValidateDefaultChannel(metabundle *MetaBundle) error {
	valid := false

	defaultChannel := metabundle.AddonMeta.DefaultChannel
	for _, channel := range metabundle.AddonMeta.Channels {
		if channel.Name == defaultChannel {
			valid = true
		}
	}

	if !valid {
		return fmt.Errorf("could not find defaultChannel in channels")
	}
	return nil
}

func GetAllMetaValidators() []Validator {
	return []Validator{
		{
			Description: "Ensure defaultChannel is present in list of channels",
			Runner:      ValidateDefaultChannel,
		},
	}
}
