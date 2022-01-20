package am0001

import (
	"context"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
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
func (d *DefaultChannel) isListedInChannels(channels *[]v1alpha1.Channel, defaultChannel string) validator.Result {
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
