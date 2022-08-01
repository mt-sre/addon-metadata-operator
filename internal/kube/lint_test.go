package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestDeploymentLinterImplInterfaces(t *testing.T) {
	t.Parallel()

	require.Implements(t, new(DeploymentLinter), new(DeploymentLinterImpl))
}

func TestDeploymentLinterImpl(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		Deployment      appsv1.Deployment
		ExpectedSuccess bool
	}{
		"default checks/valid deployment": {
			Deployment: newTestDeployment(corev1.Container{
				ReadinessProbe: &corev1.Probe{},
				LivenessProbe:  &corev1.Probe{},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewQuantity(500, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(1024, resource.BinarySI),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewQuantity(1000, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(2048, resource.BinarySI),
					},
				},
			}),
			ExpectedSuccess: true,
		},
		"default checks/missing Readiness Probe": {
			Deployment: newTestDeployment(corev1.Container{
				LivenessProbe: &corev1.Probe{},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewQuantity(500, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(1024, resource.BinarySI),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewQuantity(1000, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(2048, resource.BinarySI),
					},
				},
			}),
			ExpectedSuccess: false,
		},
		"default checks/missing liveness Probe": {
			Deployment: newTestDeployment(corev1.Container{
				ReadinessProbe: &corev1.Probe{},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewQuantity(500, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(1024, resource.BinarySI),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewQuantity(1000, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(2048, resource.BinarySI),
					},
				},
			}),
			ExpectedSuccess: false,
		},
		"default checks/missing CPU requests": {
			Deployment: newTestDeployment(corev1.Container{
				ReadinessProbe: &corev1.Probe{},
				LivenessProbe:  &corev1.Probe{},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceMemory: *resource.NewQuantity(1024, resource.BinarySI),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewQuantity(1000, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(2048, resource.BinarySI),
					},
				},
			}),
			ExpectedSuccess: false,
		},
		"default checks/missing CPU limits": {
			Deployment: newTestDeployment(corev1.Container{
				ReadinessProbe: &corev1.Probe{},
				LivenessProbe:  &corev1.Probe{},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewQuantity(500, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(1024, resource.BinarySI),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceMemory: *resource.NewQuantity(2048, resource.BinarySI),
					},
				},
			}),
			ExpectedSuccess: false,
		},
		"default checks/missing memory requests": {
			Deployment: newTestDeployment(corev1.Container{
				ReadinessProbe: &corev1.Probe{},
				LivenessProbe:  &corev1.Probe{},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU: *resource.NewQuantity(500, resource.DecimalSI),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewQuantity(1000, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(2048, resource.BinarySI),
					},
				},
			}),
			ExpectedSuccess: false,
		},
		"default checks/missing memory limits": {
			Deployment: newTestDeployment(corev1.Container{
				ReadinessProbe: &corev1.Probe{},
				LivenessProbe:  &corev1.Probe{},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewQuantity(500, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(1024, resource.BinarySI),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU: *resource.NewQuantity(1000, resource.DecimalSI),
					},
				},
			}),
			ExpectedSuccess: false,
		},
		"default checks/multiple issues": {
			Deployment:      newTestDeployment(corev1.Container{}),
			ExpectedSuccess: false,
		},
	} {
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			val := NewDeploymentLinterImpl()

			res := val.Lint(tc.Deployment)

			assert.Equal(t, tc.ExpectedSuccess, res.Success, res)
		})
	}
}

func newTestDeployment(containers ...corev1.Container) appsv1.Deployment {
	return appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: containers,
				},
			},
		},
	}
}
