package validate

import "github.com/go-playground/validator"

type Validator struct {
	Description string
	Runner      validator.StructLevelFunc
}

func ValidateDefaultChannel(sl validator.StructLevel) {
	valid := false
	metaBundle := sl.Current().Interface().(MetaBundle)

	defaultChannel := metaBundle.AddonMeta.DefaultChannel
	for _, channel := range metaBundle.AddonMeta.Channels {
		if channel.Name == defaultChannel {
			valid = true
		}
	}

	if !valid {
		sl.ReportError(metaBundle.AddonMeta.DefaultChannel,
			"defaultChannel",
			"DefaultChannel",
			"DefaultChannelPresentInChannels",
			"")
	}
}

func GetAllMetaValidators() []Validator {
	return []Validator{
		{
			Description: "Ensure defaultChannel is present in list of channels",
			Runner:      ValidateDefaultChannel,
		},
	}
}
