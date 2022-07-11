package am0007

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

func init() {
	validator.Register(NewCSVInstallModes)
}

const (
	code = 7
	name = "csv_install_modes"
	desc = "Validate installMode is supported."
)

func NewCSVInstallModes(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &CSVInstallModes{
		Base: base,
	}, nil
}

type CSVInstallModes struct {
	*validator.Base
}

func (c *CSVInstallModes) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	installMode := mb.AddonMeta.InstallMode

	// allow only AllNamespaces and OwnNamespace install mode.
	if indexOf(installMode, validInstallModes) == -1 {
		return c.Fail(fmt.Sprintf("unsupported install mode %v", installMode))
	}

	for _, bundle := range mb.Bundles {
		bundleName, err := utils.GetBundleNameVersion(bundle)
		if err != nil {
			return c.Error(err)
		}
		spec, err := extractCSVSpec(bundle)
		if err != nil {
			return c.Error(fmt.Errorf("unable to extract CSV for %v: %w", bundleName, err))
		}
		if success, failureMsg := isInstallModeSupported(spec.InstallModes, installMode); !success {
			return c.Fail(fmt.Sprintf("Bundle %v failed CSV validation: %v.", bundleName, failureMsg))
		}
	}
	return c.Success()
}

func isInstallModeSupported(installModes []operatorsv1alpha1.InstallMode, target string) (bool, string) {
	var allSupported []string

	targetSupported := false
	for _, im := range installModes {
		imType := string(im.Type)
		if im.Supported {
			allSupported = append(allSupported, imType)
		}

		if imType == target {
			targetSupported = im.Supported
		}
	}
	return targetSupported, fmt.Sprintf("Target installMode %v is not supported. CSV only supports these installModes %v.", target, allSupported)
}

type CSVSpec struct {
	InstallModes []operatorsv1alpha1.InstallMode `json:"installModes"`
}

func extractCSVSpec(b *registry.Bundle) (*CSVSpec, error) {
	csv, err := b.ClusterServiceVersion()
	if err != nil {
		return nil, err
	}
	var res CSVSpec
	if err := json.Unmarshal(csv.Spec, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

var validInstallModes = []string{"AllNamespaces", "OwnNamespace"}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}
