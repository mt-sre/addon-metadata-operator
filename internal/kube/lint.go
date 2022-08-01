package kube

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type DeploymentLinter interface {
	Lint(d appsv1.Deployment) DeploymentCheckResult
}

func NewDeploymentCheckResult(reasons ...string) DeploymentCheckResult {
	if len(reasons) > 0 {
		return DeploymentCheckResult{
			Success: false,
			Reasons: reasons,
		}
	}

	return DeploymentCheckResult{
		Success: true,
	}
}

type DeploymentCheckResult struct {
	Success bool
	Reasons []string
}

func (r DeploymentCheckResult) Join(other DeploymentCheckResult) DeploymentCheckResult {
	return DeploymentCheckResult{
		Success: r.Success && other.Success,
		Reasons: append(r.Reasons, other.Reasons...),
	}
}

type DeploymentCheckResultList []DeploymentCheckResult

func (l DeploymentCheckResultList) Reduce() DeploymentCheckResult {
	res := NewDeploymentCheckResult()

	for _, result := range l {
		res = res.Join(result)
	}

	return res
}

func NewDeploymentLinterImpl(opts ...DeploymentLinterImplOption) *DeploymentLinterImpl {
	var cfg DeploymentLinterImplConfig

	cfg.Option(opts...)
	cfg.Default()

	return &DeploymentLinterImpl{
		cfg: cfg,
	}
}

type DeploymentLinterImpl struct {
	cfg DeploymentLinterImplConfig
}

func (dv *DeploymentLinterImpl) Lint(d appsv1.Deployment) DeploymentCheckResult {
	var results DeploymentCheckResultList

	for _, check := range dv.cfg.checks {
		results = append(results, check(d))
	}

	return results.Reduce()
}

type DeploymentLinterImplConfig struct {
	checks []DeploymentCheck
}

func (c *DeploymentLinterImplConfig) Option(opts ...DeploymentLinterImplOption) {
	for _, opt := range opts {
		opt.ConfigureDeploymentValidator(c)
	}
}

func (c *DeploymentLinterImplConfig) Default() {
	if len(c.checks) == 0 {
		c.checks = []DeploymentCheck{
			HasLivenessProbes,
			HasReadinessProbes,
			HasCPUResourceRequirements,
			HasMemoryResourceRequirements,
		}
	}
}

type DeploymentLinterImplOption interface {
	ConfigureDeploymentValidator(*DeploymentLinterImplConfig)
}

type WithDeploymentChecks []DeploymentCheck

func (w WithDeploymentChecks) ConfigureDeploymentLinter(c *DeploymentLinterImplConfig) {
	c.checks = []DeploymentCheck(w)
}

type DeploymentCheck func(appsv1.Deployment) DeploymentCheckResult

func HasReadinessProbes(d appsv1.Deployment) DeploymentCheckResult {
	var reasons []string

	for _, c := range d.Spec.Template.Spec.Containers {
		if c.ReadinessProbe == nil {
			reasons = append(reasons, reportContainerLintReason(c, "missing a readiness probe"))
		}
	}

	return NewDeploymentCheckResult(reasons...)
}

func HasLivenessProbes(d appsv1.Deployment) DeploymentCheckResult {
	var reasons []string

	for _, c := range d.Spec.Template.Spec.Containers {
		if c.LivenessProbe == nil {
			reasons = append(reasons, reportContainerLintReason(c, "missing a liveness probe"))
		}
	}

	return NewDeploymentCheckResult(reasons...)
}

func HasCPUResourceRequirements(d appsv1.Deployment) DeploymentCheckResult {
	var reasons []string

	for _, c := range d.Spec.Template.Spec.Containers {
		if c.Resources.Requests.Cpu().IsZero() {
			reasons = append(reasons, reportContainerLintReason(c, "missing CPU requests"))
		}

		if c.Resources.Limits.Cpu().IsZero() {
			reasons = append(reasons, reportContainerLintReason(c, "missing CPU limits"))
		}
	}

	return NewDeploymentCheckResult(reasons...)
}

func HasMemoryResourceRequirements(d appsv1.Deployment) DeploymentCheckResult {
	var reasons []string

	for _, c := range d.Spec.Template.Spec.Containers {
		if c.Resources.Requests.Memory().IsZero() {
			reasons = append(reasons, reportContainerLintReason(c, "missing memory requests"))
		}

		if c.Resources.Limits.Memory().IsZero() {
			reasons = append(reasons, reportContainerLintReason(c, "missing memory limits"))
		}
	}

	return NewDeploymentCheckResult(reasons...)
}

func reportContainerLintReason(c corev1.Container, reason string) string {
	return fmt.Sprintf("container %q is %s", c.Name, reason)
}
