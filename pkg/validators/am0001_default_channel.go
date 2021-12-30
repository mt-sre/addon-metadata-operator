package validators

import (
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

func init() {
	Registry.Add(AM0001)
}

var AM0001 = types.Validator{
	Code:        "AM0001",
	Name:        "default_channel",
	Description: "Ensure defaultChannel is present in list of channels",
	Runner:      validateDefaultChannel,
}

// validateDefaultChannel validates whether the 'defaultChannel' provided under an addon.yaml also exists under 'channels' field
func validateDefaultChannel(mb types.MetaBundle) types.ValidatorResult {
	defaultChannel := mb.AddonMeta.DefaultChannel
	channels := mb.AddonMeta.Channels

	if res := isPartOfEnum(defaultChannel); !res.IsSuccess() {
		return res
	}
	// to be deprecated - only used for legacy builds
	if res := isListedInChannels(channels, defaultChannel); !res.IsSuccess() {
		return res
	}

	return Success()
}

func isPartOfEnum(defaultChannel string) types.ValidatorResult {
	enum := map[string]struct{}{
		"alpha":  {},
		"beta":   {},
		"stable": {},
		"edge":   {},
		"rc":     {},
	}
	if _, ok := enum[defaultChannel]; !ok {
		msg := fmt.Sprintf("The defaultChannel '%v' is not part of the accepted values: alpha, beta, stable, edge or rc.", defaultChannel)
		return Fail(msg)
	}
	return Success()
}

// TODO - deprecate this when we remove legacy builds
func isListedInChannels(channels *[]v1alpha1.Channel, defaultChannel string) types.ValidatorResult {
	// as the Channels field is deprecated, it can be omitted
	if channels == nil {
		return Success()
	}
	var channelNames []string
	for _, channel := range *channels {
		if channel.Name == defaultChannel {
			return Success()
		}
		channelNames = append(channelNames, channel.Name)
	}
	msg := fmt.Sprintf("The defaultChannel '%v' is not part of the listed channelNames: %v.", defaultChannel, channelNames)
	return Fail(msg)
}
