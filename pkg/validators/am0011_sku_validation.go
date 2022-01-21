package validators

import (
	"context"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

var AM0011 = types.Validator{
	Code:        "AM0011",
	Name:        "sku_validation",
	Description: "Validates whether a SKU Rule exists in OCM for quota provided in addon metadata",
	Runner:      ValidateOCMSKUExists,
}

func init() {
	Registry.Add(AM0011)
}

func ValidateOCMSKUExists(mb types.MetaBundle) types.ValidatorResult {
	client, err := utils.NewDefaultOCMClient()
	if err != nil {
		return Error(err)
	}

	validateFunc := GenerateOCMSKUValidator(client)

	return validateFunc(mb)
}

func GenerateOCMSKUValidator(ocm OCMClient) types.ValidateFunc {
	return func(mb types.MetaBundle) types.ValidatorResult {
		// addons with '0' quota cost are not processed for SKU validation
		if mb.AddonMeta.OcmQuotaCost == 0 {
			return Success()
		}

		// Will become the caller's responsibility to provide this in the future
		ctx := context.Background()

		quotaName := mb.AddonMeta.OcmQuotaName

		skuRules, err := ocm.GetSKURules(ctx, quotaName)
		if err != nil {
			if IsOCMServerSideError(err) {
				return RetryableError(err)
			}

			return Error(err)
		}

		if len(skuRules) < 1 {
			return Fail(fmt.Sprintf("no SKU Rule exists for ocmQuotaName '%s'", quotaName))
		}

		if err := ocm.CloseConnection(); err != nil {
			return Error(err)
		}

		return Success()
	}
}
