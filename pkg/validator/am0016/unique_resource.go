package am0016

import (
	"context"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewUniqueResource)
}

const (
	code = 16
	name = "unique_resource"
	desc = "Ensure that addon additional catalog source, secrets and credential requests names are unique"
)

func NewUniqueResource(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &UniqueResource{
		Base: base,
	}, nil
}

type UniqueResource struct {
	*validator.Base
}

func (u *UniqueResource) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	var messages []string

	acsMessage := UniqueAdditionalCataloSources(mb.AddonMeta)
	if len(acsMessage) > 0 {
		messages = append(messages, acsMessage...)
	}

	secretMessage := UniqueSecrets(mb.AddonMeta)
	if len(secretMessage) > 0 {
		messages = append(messages, secretMessage...)
	}

	if len(messages) > 0 {
		return u.Fail(messages...)
	}
	return u.Success()
}

func UniqueAdditionalCataloSources(addonMetadataSpec *v1alpha1.AddonMetadataSpec) []string {
	var messages []string

	additionalCatalogSources := addonMetadataSpec.AdditionalCatalogSources
	if additionalCatalogSources == nil {
		return messages
	}
	if len(*additionalCatalogSources) == 1 {
		return messages
	}

	catalogSourceNames := make(map[string]bool)
	for _, additionalCatalogSource := range *additionalCatalogSources {
		catalogSourceName := additionalCatalogSource.Name

		if catalogSourceNames[catalogSourceName] {
			messages = append(messages,
				fmt.Sprintf("additional catalaog source: additionalCatalogSource name %v is already present and not unique.", catalogSourceName),
			)
		} else {
			catalogSourceNames[catalogSourceName] = true
		}
	}
	return messages
}

func UniqueSecrets(addonMetadataSpec *v1alpha1.AddonMetadataSpec) []string {
	var messages []string

	config := addonMetadataSpec.Config
	// if config is nil
	if config == nil {
		return messages
	}

	secrets := config.Secrets
	if secrets == nil {
		return messages
	}
	if len(*secrets) == 1 {
		return messages
	}

	secretNames := make(map[string]bool)
	for _, secret := range *secrets {
		secretName := secret.Name

		if secretNames[secretName] {
			messages = append(messages, fmt.Sprintf("secrets: secret name %v is already present and not unique.", secretName))
		} else {
			secretNames[secretName] = true
		}
	}

	return messages
}
