package am0003

import (
	"context"
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

const (
	code = 3
	name = "operator_name"
	desc = "Validate the operatorName matches csv.Name, csv.Replaces and bundle package annotation."
)

func init() {
	validator.Register(NewOperatorName)
}

func NewOperatorName(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}
	return &OperatorName{
		Base: base,
	}, nil
}

type OperatorName struct {
	*validator.Base
}

func (o *OperatorName) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	operatorName, bundles := mb.AddonMeta.OperatorName, mb.Bundles

	var failures []string

	for _, b := range bundles {
		msgs, err := validateBundle(b, operatorName)
		if err != nil {
			return o.Error(err)
		}

		failures = append(failures, msgs...)
	}

	if len(failures) > 0 {
		return o.Fail(failures...)
	}

	return o.Success()
}

func validateBundle(bundle *registry.Bundle, operatorName string) ([]string, error) {
	var msgs []string

	nameVer, err := utils.GetBundleNameVersion(bundle)
	if err != nil {
		return nil, fmt.Errorf("retrieving bundle name and version: %w", err)
	}

	if bundle.Annotations.PackageName != operatorName {
		msgs = append(msgs, fmt.Sprintf(
			"bundle %q package annotation does not match operatorName %q", nameVer, operatorName,
		))
	}

	csv, err := bundle.ClusterServiceVersion()
	if err != nil {
		return nil, fmt.Errorf("retrieving csv: %w", err)
	}

	if msg := validateBundleIdentifier(csv.Name, operatorName); msg != "" {
		msgs = append(msgs, fmt.Sprintf("bundle %q failed validation on csv.Name: %s.", nameVer, msg))
	}

	replaces, err := csv.GetReplaces()
	if err != nil {
		return nil, fmt.Errorf("retrieving 'replaces' field from csv: %w", err)
	}

	if replaces != "" {
		if msg := validateBundleIdentifier(replaces, operatorName); msg != "" {
			msgs = append(msgs, fmt.Sprintf("bundle %q failed validation on csv.Replaces: %s.", nameVer, msg))
		}
	}

	return msgs, nil
}

func validateBundleIdentifier(id, operatorName string) string {
	parts := strings.SplitN(id, ".", 2)

	if len(parts) != 2 {
		return fmt.Sprintf("invalid format %q; expected '<name>.<version>'.", id)
	}

	bundleOperatorName, bundleVersion := parts[0], parts[1]

	if bundleOperatorName != operatorName {
		return fmt.Sprintf("invalid operatorName for %q; expected %q", id, operatorName)
	}

	if _, err := semver.ParseTolerant(bundleVersion); err != nil {
		return fmt.Sprintf("invalid semver %q: %v", id, err)
	}

	return ""
}
