package validators

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

func init() {
	Registry.Add(AM0004)
}

var AM0004 = types.NewValidator(
	"AM0004",
	types.ValidateFunc(ValidateIconBase64),
	types.ValidatorName("icon_base64"),
	types.ValidatorDescription("Ensure that `icon` in Addon metadata is rightfully base64 encoded"),
)

// ValidateIconBase64 validates 'icon' in the addon metadata is rightfully base64 encoded
func ValidateIconBase64(cfg types.ValidatorConfig, mb types.MetaBundle) types.ValidatorResult {
	icon := mb.AddonMeta.Icon
	if icon == "" {
		return Fail(fmt.Sprintf("`icon` not found under the addon metadata of %s", mb.AddonMeta.ID))
	}

	b64decoded, err := base64.StdEncoding.DecodeString(icon)
	if err != nil {
		return Fail(fmt.Sprintf("`icon` found to be improperly base64 populated under the addon metadata of %s", mb.AddonMeta.ID))
	}

	_, err = png.Decode(bytes.NewReader(b64decoded))
	if err != nil {
		return Fail(fmt.Sprintf("`icon`'s base64 value found to correspond to a non-png data under the addon metadata of %s", mb.AddonMeta.ID))
	}

	return Success()
}
