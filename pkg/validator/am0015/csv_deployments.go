package am0015

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	operatorv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/registry"
	"golang.stackrox.io/kube-linter/pkg/checkregistry"
	"golang.stackrox.io/kube-linter/pkg/lintcontext"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
		Base: base,
	}, nil
}

type CSVDeployment struct {
	*validator.Base
}

type Spec struct {
	InstallStrategy operatorv1alpha1.NamedInstallStrategy `json:"install"`
}

func (c *CSVDeployment) Run(ctx context.Context, mb types.MetaBundle) validator.Result {

	// list of checks to be performed under CSV Deployment validation
	checks := message{
		"no-liveness-probe",
		"no-readiness-probe",
		"unset-cpu-requirements",
		"unset-memory-requirements",
	}

	if res := c.isValidCSVDeployment(mb.Bundles, checks); !res.IsSuccess() {
		return res
	}
	return c.Success()
}

type message []string

func (c *CSVDeployment) isValidCSVDeployment(bundles []*registry.Bundle, checks message) validator.Result {
	var msg message
	var registry checkregistry.CheckRegistry
	// to check whether the check is present or not
	instantiatedcheck, err := GetInstantiatedChecks(registry, msg)
	if err != nil {
		msg = append(msg, fmt.Sprintf("Checks are not present in the registry"))
		return c.Fail(msg...)
	}

	for _, bundle := range bundles {
		csv, err := bundle.ClusterServiceVersion()
		if err != nil {
			msg = append(msg, fmt.Sprintf("error %v occured for extracting CSV for bundle %v", err, bundle.Name))
			continue
		}

		var spec Spec
		if err := json.Unmarshal(csv.Spec, &spec); err != nil {
			return c.Error(err)
		}

		var objs []client.Object
		for _, d := range spec.InstallStrategy.StrategySpec.DeploymentSpecs {
			objs = append(objs, &appsv1.Deployment{
				Spec: d.Spec,
			})
		}

		lintCtxs := []lintcontext.LintContext{}
		lintCtx := &validator.LintContextImpl{}

		lintCtx.AddObjects(lintcontext.Object{K8sObject: objs})
		lintCtxs = append(lintCtxs, lintCtx)
		result := RunValidations(lintCtxs, instantiatedcheck)
		if result == ChecksFailed {
			msg = append(msg, fmt.Sprintf("checks failed for bundle %v", bundle.Name))
		}
	}

	if len(msg) > 0 {
		return c.Fail(msg...)
	}
	return c.Success()
}
