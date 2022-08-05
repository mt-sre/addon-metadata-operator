package am0015

import (
	"context"
	"encoding/json"

	"github.com/mt-sre/addon-metadata-operator/internal/kube"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	operatorv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
)

func init() {
	validator.Register(NewCSVDeployment)
}

const (
	code = 15
	name = "csv_deployments"
	desc = "Ensure all deployment in CSV must have valid resource requests, livenessprobe and readinessprobe"
)

func NewCSVDeployment(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &CSVDeployment{
		Base:   base,
		linter: kube.NewDeploymentLinterImpl(),
	}, nil
}

type CSVDeployment struct {
	*validator.Base
	linter kube.DeploymentLinter
}

type Spec struct {
	InstallStrategy operatorv1alpha1.NamedInstallStrategy `json:"install"`
}

func (c *CSVDeployment) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	var msgs []string
	var spec Spec
	bundle, err := validator.GetLatestBundle(mb.Bundles)
	if err != nil {
		c.Fail("Error while checking bundles")
	}

	csv, err := bundle.ClusterServiceVersion()
	if err != nil {
		c.Error(err)
	}

	if err := json.Unmarshal(csv.Spec, &spec); err != nil {
		c.Error(err)
	}

	for _, deploymentSpec := range spec.InstallStrategy.StrategySpec.DeploymentSpecs {
		deployment := appsv1.Deployment{Spec: deploymentSpec.Spec}
		res := c.linter.Lint(deployment)
		if !res.Success {
			msgs = append(msgs, res.Reasons...)
		}
	}

	if len(msgs) > 0 {
		return c.Fail(msgs...)
	}
	return c.Success()
}
