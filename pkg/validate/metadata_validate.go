package validate

func ValidateDefaultChannel(metabundle *MetaBundle) (bool, error) {
	defaultChannel := metabundle.AddonMeta.DefaultChannel
	for _, channel := range metabundle.AddonMeta.Channels {
		if channel.Name == defaultChannel {
			return true, nil
		}
	}
	return false, nil
}

func GetAllMetaValidators() []Validator {
	return []Validator{
		{
			Description: "Ensure defaultChannel is present in list of channels",
			Runner:      ValidateDefaultChannel,
		},
	}
}
