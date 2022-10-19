package am0017

import (
	"context"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewSecret)
}

const (
	code = 17
	name = "secrets"
	desc = "Ensure that addon secrets names are unique"
)

func NewSecret(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}
	return &Secret{
		Base: base,
	}, nil
}

type Secret struct {
	*validator.Base
}

func (s *Secret) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	config := mb.AddonMeta.Config
	// if config is nil
	if config == nil {
		return s.Success()
	}

	secrets := config.Secrets
	if secrets == nil {
		return s.Success()
	}
	if len(*secrets) == 1 {
		return s.Success()
	}

	var messages []string
	secretNames := make(map[string]bool)
	for _, secret := range *secrets {
		secretName := secret.Name

		if secretNames[secretName] {
			messages = append(messages, fmt.Sprintf("secret name %v is already present and not unique.", secretName))
		} else {
			secretNames[secretName] = true
		}
	}

	if len(messages) > 0 {
		return s.Fail(messages...)
	}
	return s.Success()
}
