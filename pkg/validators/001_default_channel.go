package validators

import (
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

// Validate001DefaultChannel validates whether the 'defaultChannel' provided under an addon.yaml also exists under 'channels' field
func Validate001DefaultChannel(metabundle utils.MetaBundle) (bool, string, error) {
	defaultChannel := metabundle.AddonMeta.DefaultChannel
	channels := metabundle.AddonMeta.Channels

	if success, failureMsg := isPartOfEnum(defaultChannel); !success {
		return false, failureMsg, nil
	}
	// to be deprecated - only used for legacy builds
	if success, failureMsg := isListedInChannels(channels, defaultChannel); !success {
		return false, failureMsg, nil
	}

	// if success, failureMsg, err := matchesDefaultChannelAnnotations(metabundle.Bundles, defaultChannel); !success {
	// 	return false, failureMsg, err
	// }
	return true, "", nil
}

func isPartOfEnum(defaultChannel string) (bool, string) {
	enum := map[string]struct{}{
		"alpha":  {},
		"beta":   {},
		"stable": {},
		"edge":   {},
		"rc":     {},
	}
	if _, ok := enum[defaultChannel]; !ok {
		return false, fmt.Sprintf("The defaultChannel '%v' is not part of the accepted values: alpha, beta, stable, edge or rc.", defaultChannel)
	}
	return true, ""
}

// TODO - deprecate this when we remove legacy builds
func isListedInChannels(channels *[]v1alpha1.Channel, defaultChannel string) (bool, string) {
	// as the Channels field is deprecated, it can be omitted
	if channels == nil {
		return true, ""
	}
	var channelNames []string
	for _, channel := range *channels {
		if channel.Name == defaultChannel {
			return true, ""
		}
		channelNames = append(channelNames, channel.Name)
	}
	return false, fmt.Sprintf("The defaultChannel '%v' is not part of the listed channelNames: %v.", defaultChannel, channelNames)
}

// TODO - (sblaisdo) enable when extraction format is bundles instead of packageManifest
// if the annotation is not present, make sure the defaultChannel is alpha
// func matchesDefaultChannelAnnotations(bundles []registry.Bundle, defaultChannel string) (bool, string, error) {
// 	for _, bundle := range bundles {
// 		bundleName, err := utils.GetBundleNameVersion()
// 		if err != nil {
// 			return false, "", err
// 		}
// 		failureMsgs := []string{fmt.Sprintf("Missing or invalid channel annotation for bundle '%v'", bundleName)}

// 		if bundle.Annotations == nil {
// 			err := fmt.Errorf("bundles.Annotations is nil for '%v'. The extractor should have reported an error.", bundleName)
// 			return false, strings.Join(failureMsgs, ": "), err
// 		}

// 		allChannels := strings.Split(bundle.Annotations.Channels, ",")
// 		var defaultBundleChannel string
// 		if len(allChannels) == 1 {
// 			// if a single channel is listed, should match defaultChannel
// 			defaultBundleChannel = allChannels[0]
// 			failureMsgs = append(failureMsgs, fmt.Sprintf("Please set the 'operators.operatorframework.io.bundle.channels.v1' annotation to '%v'.", defaultChannel))

// 		} else {
// 			// if multiple channels are listed, need to specify the optional
// 			// operators.operatorframework.io.bundle.channel.default.v1
// 			defaultBundleChannel = bundle.Annotations.DefaultChannelName
// 			failureMsgs = append(failureMsgs, fmt.Sprintf("Please set the operators.operatorframework.io.bundle.channel.default.v1 annotation to '%v'.", defaultChannel))
// 		}
// 		if defaultChannel != defaultBundleChannel {
// 			failureMsgs = append(failureMsgs, fmt.Sprintf("DefaultChannel '%v' does not match channel from bundle annotations '%v'.", defaultChannel, defaultBundleChannel))
// 			return false, strings.Join(failureMsgs, ": "), nil
// 		}
// 	}
// 	return true, "", nil
// }
