package am0017

import (
	"context"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewPullSecretname)
}

const (
	code = 18
	name = "pull_secret_name"
	desc = "Ensure that pullSecretName if not nil is present in Secrets"
)

func NewPullSecretname(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}
	return &PullSecretname{
		Base: base,
	}, nil
}

type PullSecretname struct {
	*validator.Base
}

func (p *PullSecretname) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	pullSecretName := mb.AddonMeta.PullSecretName
	if pullSecretName == "" {
		return p.Success()
	}

	config := mb.AddonMeta.Config
	// if config is nil
	if config == nil {
		return p.Fail(fmt.Sprintf("pullSecretName %v is present in addon.yaml whereas addon config is nil", pullSecretName))
	}

	secrets := config.Secrets
	// if secrets is nil
	if secrets == nil {
		return p.Fail(fmt.Sprintf("pullSecretName %v is present in addon.yaml whereas addon secrets are nil", pullSecretName))
	}

	for _, secret := range *secrets {
		if pullSecretName == secret.Name {
			return p.Success()
		}
	}
	return p.Fail(fmt.Sprintf("pullSecretName %v is not present in addon secrets", pullSecretName))
}
