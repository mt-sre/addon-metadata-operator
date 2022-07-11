package types

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

type MetaBundle struct {
	AddonMeta *v1alpha1.AddonMetadataSpec
	Bundles   []*registry.Bundle
}

func NewMetaBundle(addonMeta *v1alpha1.AddonMetadataSpec, bundles []*registry.Bundle) *MetaBundle {
	return &MetaBundle{
		AddonMeta: addonMeta,
		Bundles:   bundles,
	}
}
