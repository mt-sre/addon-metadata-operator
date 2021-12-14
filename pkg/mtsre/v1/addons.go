package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
This package is for the types/sections of addon metadata schema which aren't compliant with OCM API Spec.
Please keep in sync with managed-tenants-cli schemas:
- https://github.com/mt-sre/managed-tenants-cli/blob/main/managedtenants/data/metadata.schema.yaml
- https://github.com/mt-sre/managed-tenants-cli/blob/main/managedtenants/data/imageset.schema.yaml

We need to generate deepcopy methods for all complex types. Using the non-root
generation annotation as these types don't implement runtime.Object interface.

Update zz_generated.deepcopy.go with:
	$ make generate
*/

// +kubebuilder:validation:Pattern=`^([A-Za-z -]+ <[0-9A-Za-z_.-]+@redhat\.com>,?)+$`
type Notification string

//+kubebuilder:object:generate=true
type Monitoring struct {
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace" validate:"required"`

	// +kubebuilder:validation:Required
	MatchNames []string `json:"matchNames" validate:"required"`

	// +kubebuilder:validation:Required
	MatchLabels map[string]string `json:"matchLabels" validate:"required"`
}

//+kubebuilder:object:generate=true
type BundleParameters struct {
	// +optional
	// +kubebuilder:validation:Pattern=`^(true|false|^$)$`
	UseClusterStorage *string `json:"useClusterStorage"`

	// +optional
	// +kubebuilder:validation:Pattern=`^([0-9A-Za-z_.-]+@redhat\.com,? ?)+$`
	AlertingEmailAddress *string `json:"alertingEmailAddress"`

	// +optional
	// +kubebuilder:validation:Pattern=`^([0-9A-Za-z_.-]+@redhat\.com,? ?)+$`
	BuAlertingEmailAddress *string `json:"buAlertingEmailAddress"`

	// +optional
	// +kubebuilder:validation:Pattern=`^[0-9A-Za-z._-]+@(devshift\.net|rhmw\.io)$`
	AlertSMTPFrom *string `json:"alertSMTPFrom"`

	// +optional
	// +kubebuilder:validation:Pattern=`^addon-[0-9A-Za-z-]+-parameters$`
	AddonParamsSecretName *string `json:"addonParamsSecretName"`
}

//+kubebuilder:object:generate=true
type PagerDuty struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9]+$`
	EscalationPolicy string `json:"snitchNamePostFix" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	AcknowledgeTimeout int `json:"acknowledgeTimeout" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	ResolveTimeout int `json:"resolveTimeout" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	SecretName string `json:"secretName" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	SecretNamespace string `json:"secretNamespace" validate:"required"`
}

//+kubebuilder:object:generate=true
type DeadmansSnitch struct {
	// +optional
	ClusterDeploymentSelector *metav1.LabelSelector `json:"clusterDeploymentSelector"`

	// +optional
	SnitchNamePostFix *string `json:"snitchNamePostFix"`

	// +optional
	TargetSecretRef *TargetSecretRef `json:"targetSecretRef"`

	// +kubebuilder:validation:Required
	Tags []Tag `json:"tags" validate:"required"`
}

// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
type Tag string

//+kubebuilder:object:generate=true
type TargetSecretRef struct {
	// +optional
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	Name *string `json:"name"`

	// +optional
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	Namespace *string `json:"namespace"`
}

//+kubebuilder:object:generate=true
type SubscriptionConfig struct {
	// +kubebuilder:validation:Required
	Env *[]EnvItem `json:"env" validate:"required"`
}

//+kubebuilder:object:generate=true
type EnvItem struct {
	// +kubebuilder:validation:Required
	Name string `json:"name" validate:"required"`

	// +kubebuilder:validation:Required
	Value string `json:"value" validate:"required"`
}
