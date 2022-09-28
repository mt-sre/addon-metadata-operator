package am0015

import (
	"context"

	"github.com/mt-sre/addon-metadata-operator/internal/kube"
	"github.com/mt-sre/addon-metadata-operator/pkg/operator"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
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

func (c *CSVDeployment) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	var msgs []string

	bundle, ok := operator.HeadBundle(mb.Bundles...)
	if !ok {
		return c.Success()
	}

	csv := bundle.ClusterServiceVersion

	for _, deploymentSpec := range csv.Spec.InstallStrategy.StrategySpec.DeploymentSpecs {
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
