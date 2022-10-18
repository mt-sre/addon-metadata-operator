package am0016

import (
	"context"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewAdditionalCatalogSource)
}

const (
	code = 16
	name = "additional_catalog_sources"
	desc = "Ensure that addon additional catalog source names must be unique"
)

func NewAdditionalCatalogSource(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &AdditionalCatalogSource{
		Base: base,
	}, nil
}

type AdditionalCatalogSource struct {
	*validator.Base
}

func (a *AdditionalCatalogSource) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	additionalCatalogSources := mb.AddonMeta.AdditionalCatalogSources

	if additionalCatalogSources == nil {
		return a.Success()
	}
	if len(*additionalCatalogSources) == 1 {
		return a.Success()
	}

	var messages []string
	catalogSourceNames := make(map[string]bool)
	for _, additionalCatalogSource := range *additionalCatalogSources {
		catalogSourceName := additionalCatalogSource.Name

		if catalogSourceNames[catalogSourceName] {
			messages = append(messages, fmt.Sprintf("additionalCatalogSource name %v is already present and not unique.", catalogSourceName))
		} else {
			catalogSourceNames[catalogSourceName] = true
		}
	}
	if len(messages) > 0 {
		return a.Fail(messages...)
	}
	return a.Success()
}
