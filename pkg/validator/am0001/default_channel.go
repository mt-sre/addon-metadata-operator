package am0001

import (
	"context"
	"fmt"

	opsv1alpha1 "github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/operator"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewDefaultChannel)
}

const (
	code = 1
	name = "default_channel"
	desc = "Ensure defaultChannel is present in list of channels"
)

func NewDefaultChannel(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &DefaultChannel{
		Base: base,
	}, nil
}

type DefaultChannel struct {
	*validator.Base
}

func (d *DefaultChannel) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	defaultChannel := mb.AddonMeta.DefaultChannel
	channels := mb.AddonMeta.Channels

	if res := d.isPartOfEnum(defaultChannel); !res.IsSuccess() {
		return res
	}
	// to be deprecated - only used for legacy builds
	if res := d.isListedInChannels(channels, defaultChannel); !res.IsSuccess() {
		return res
	}

	if res := d.matchesBundleChannelAnnotations(defaultChannel, mb.Bundles); !res.IsSuccess() {
		return res
	}

	return d.Success()
}

func (d *DefaultChannel) isPartOfEnum(defaultChannel string) validator.Result {
	enum := map[string]struct{}{
		"alpha":  {},
		"beta":   {},
		"stable": {},
		"edge":   {},
		"rc":     {},
	}
	if _, ok := enum[defaultChannel]; !ok {
		msg := fmt.Sprintf("The defaultChannel '%v' is not part of the accepted values: alpha, beta, stable, edge or rc.", defaultChannel)
		return d.Fail(msg)
	}
	return d.Success()
}

// TODO - deprecate this when we remove legacy builds
func (d *DefaultChannel) isListedInChannels(channels *[]opsv1alpha1.Channel, defaultChannel string) validator.Result {
	// as the Channels field is deprecated, it can be omitted
	if channels == nil {
		return d.Success()
	}
	var channelNames []string
	for _, channel := range *channels {
		if channel.Name == defaultChannel {
			return d.Success()
		}
		channelNames = append(channelNames, channel.Name)
	}
	msg := fmt.Sprintf("The defaultChannel '%v' is not part of the listed channelNames: %v.", defaultChannel, channelNames)
	return d.Fail(msg)
}

func (d *DefaultChannel) matchesBundleChannelAnnotations(defaultChannel string, bundles []operator.Bundle) validator.Result {
	var message []string

	bundle, ok := operator.HeadBundle(bundles...)
	if !ok {
		return d.Success()
	}

	if bundle.Annotations.DefaultChannelName == "" && defaultChannel != "alpha" {
		msg := fmt.Sprintf("operators.operatorframework.io.bundle.channel.default.v1 is not defined so defaultChannel should be 'alpha' instead of '%v'", defaultChannel)
		message = append(message, msg)
	}

	if bundle.Annotations.DefaultChannelName != defaultChannel && bundle.Annotations.DefaultChannelName != "" {
		msg := fmt.Sprintf("The defaultChannel '%v' does not match annotation operators.operatorframework.io.bundle.channel.default.v1 '%v'.",
			defaultChannel, bundle.Annotations.DefaultChannelName,
		)
		message = append(message, msg)
	}

	channels := bundle.Annotations.Channels

	if !isPresentInBundleChannels(defaultChannel, channels) {
		msg := fmt.Sprintf("The defaultChannel '%v' is not present in annotation operators.operatorframework.io.bundle.channels.v1 '%v'.",
			defaultChannel, channels,
		)
		message = append(message, msg)
	}

	if len(message) > 0 {
		return d.Fail(message...)
	}
	return d.Success()
}

func isPresentInBundleChannels(defaultChannel string, channels []string) bool {
	for _, channel := range channels {
		if channel == defaultChannel {
			return true
		}
	}
	return false
}
