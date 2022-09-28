package am0012

import (
	"context"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/operator"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils/csvutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewCSVRBAC)
}

const (
	code = 12
	name = "csv_permissions"
	desc = "Validates the permissions specified in the csv"
)

func NewCSVRBAC(opt validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &CSVRBAC{
		Base: base,
	}, nil
}

type CSVRBAC struct {
	*validator.Base
}

func (v *CSVRBAC) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	bundle, ok := operator.HeadBundle(mb.Bundles...)
	if !ok {
		return v.Success()
	}

	csv := bundle.ClusterServiceVersion

	apisOwnedByOperator, err := csvutils.GetApisOwned(csv)
	if err != nil {
		return v.Error(err)
	}

	permissions, err := csvutils.GetPermissions(csv)
	if err != nil {
		return v.Error(err)
	}

	validationErrors := []string{}
	validationErrors, err = validateApiGroups(permissions, validationErrors)
	if err != nil {
		return v.Error(err)
	}
	validationErrors, err = validateWildcardInResources(permissions, apisOwnedByOperator, validationErrors)
	if err != nil {
		return v.Error(err)
	}
	validationErrors, err = validateConfidentialObjAccessAtClusterScope(permissions, validationErrors)
	if err != nil {
		return v.Error(err)
	}

	if len(validationErrors) == 0 {
		return v.Success()
	}
	failureMsg := "CSV rbac validation errors: \n" + strings.Join(validationErrors, "\n")
	return v.Fail(failureMsg)
}

func validateApiGroups(permissions *types.CSVPermissions, existingValidationErrs []string) ([]string, error) {
	if csvutils.WildCardApiGroupPresent(permissions) {
		errorMsg := "Wild card string used under api group/s"
		return append(existingValidationErrs, errorMsg), nil
	}
	return existingValidationErrs, nil
}

func validateWildcardInResources(csvPermissions *types.CSVPermissions,
	apisOwnedByOperator []string,
	existingValidationErrs []string) ([]string, error) {
	if csvutils.WildCardResourcePresent(csvPermissions, apisOwnedByOperator) {
		errorMsg := "Wild card string used under resource/s not owned by the operator"
		return append(existingValidationErrs, errorMsg), nil
	}
	return existingValidationErrs, nil
}

func validateConfidentialObjAccessAtClusterScope(csvPermissions *types.CSVPermissions,
	existingValidationErrs []string) ([]string, error) {
	if csvutils.CheckForConfidentialObjAccessAtClusterScope(csvPermissions) {
		errorMsg := "config maps/Secrets access rules present at the cluster scope"
		return append(existingValidationErrs, errorMsg), nil
	}
	return existingValidationErrs, nil
}
