package types

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	op "github.com/mt-sre/addon-metadata-operator/pkg/operator"
)

type MetaBundle struct {
	AddonMeta *v1alpha1.AddonMetadataSpec
	Bundles   []op.Bundle
}

func NewMetaBundle(addonMeta *v1alpha1.AddonMetadataSpec, bundles []op.Bundle) *MetaBundle {
	return &MetaBundle{
		AddonMeta: addonMeta,
		Bundles:   bundles,
	}
}
