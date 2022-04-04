package am0011

import (
	"context"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewOCMSKURuleExists)
}

const (
	code = 11
	name = "sku_validation"
	desc = "Validates whether a SKU Rule exists in OCM for quota provided in addon metadata"
)

func NewOCMSKURuleExists(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &OCMSKURuleExists{
		Base: base,
		ocm:  deps.OCMClient,
	}, nil
}

type OCMSKURuleExists struct {
	*validator.Base
	ocm validator.QuotaRuleGetter
}

func (o *OCMSKURuleExists) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	quotaName := mb.AddonMeta.OcmQuotaName

	quotaRuleExists, err := o.ocm.QuotaRuleExists(ctx, quotaName)
	if err != nil {
		if validator.IsOCMServerSideError(err) {
			return o.RetryableError(err)
		}

		return o.Error(err)
	}

	if !quotaRuleExists {
		return o.Fail(fmt.Sprintf("no QuotaRule exists for ocmQuotaName '%s'", quotaName))
	}

	return o.Success()
}
