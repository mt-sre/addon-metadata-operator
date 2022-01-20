package am0004

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewIconBase64)
}

const (
	code = 4
	name = "icon_base64"
	desc = "Ensure that `icon` in Addon metadata is rightfully base64 encoded"
)

func NewIconBase64(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &IconBase64{
		Base: base,
	}, nil
}

type IconBase64 struct {
	*validator.Base
}

// ValidateIconBase64 validates 'icon' in the addon metadata is rightfully base64 encoded
func (i *IconBase64) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	icon := mb.AddonMeta.Icon
	if icon == "" {
		return i.Fail(fmt.Sprintf("`icon` not found under the addon metadata of %s", mb.AddonMeta.ID))
	}

	b64decoded, err := base64.StdEncoding.DecodeString(icon)
	if err != nil {
		return i.Fail(fmt.Sprintf("`icon` found to be improperly base64 populated under the addon metadata of %s", mb.AddonMeta.ID))
	}

	_, err = png.Decode(bytes.NewReader(b64decoded))
	if err != nil {
		return i.Fail(fmt.Sprintf("`icon`'s base64 value found to correspond to a non-png data under the addon metadata of %s", mb.AddonMeta.ID))
	}

	return i.Success()
}
