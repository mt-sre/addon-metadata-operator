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

// Deprecated: Replaced by MetricsFederation
// +kubebuilder:object:generate=true
type Monitoring struct {
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace" validate:"required"`

	// +kubebuilder:validation:Required
	MatchNames []string `json:"matchNames" validate:"required"`

	// +kubebuilder:validation:Required
	MatchLabels map[string]string `json:"matchLabels" validate:"required"`
}

// +kubebuilder:object:generate=true
type MetricsFederation struct {
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	Namespace string `json:"namespace" validate:"required"`

	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	PortName string `json:"portName" validate:"required"`

	// +kubebuilder:validation:Pattern=`^[a-zA-Z_:][a-zA-Z0-9_:]*$`
	MatchNames []string `json:"matchNames" validate:"required"`

	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9-_./]+$`
	MatchLabels map[string]string `json:"matchLabels" validate:"required"`
}

// +kubebuilder:object:generate=true
type MonitoringStack struct {
	// +optional
	// This denotes whether the addon requires the MonitoringStack CR to be created in runtime or not. Validation fails if it is provided as 'false' and at the same time other parameters are specified
	Enabled *bool `json:"enabled,omitempty"`

	// +optional
	// Represents the resource quotas (requests/limits) to be allocated to the Prometheus instances which will be spun up consequently by the respective MonitoringStack CR in runtime. If not provided, the default values would be used: '{requests: {cpu: '100m', memory: '256M'}, limits:{memory: '512M', cpu: '500m'}}'
	Resources *MonitoringStackResources `json:"resources,omitempty"`
}

type MonitoringStackResources struct {
	// Represents the cpu/memory resources which would be requested by the Prometheus instances spun up consequently by the MonitoringStack CR in runtime
	Request *MonitoringStackResource `json:"requests,omitempty"`

	// Represents the max. amount of cpu/memory resources which would be accessible by the Prometheus instances spun up consequently by the MonitoringStack CR in runtime
	Limits *MonitoringStackResource `json:"limits,omitempty"`
}

type MonitoringStackResource struct {
	// Ref: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go#L147
	// +kubebuilder:validation:Pattern=`^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$`
	Cpu *string `json:"cpu,omitempty"`

	// +kubebuilder:validation:Pattern=`^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$`
	Memory *string `json:"memory,omitempty"`
}

// Deprecated: Replaced by SubscriptionConfig.
// +kubebuilder:object:generate=true
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

// +kubebuilder:object:generate=true
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

// +kubebuilder:object:generate=true
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

// +kubebuilder:object:generate=true
type TargetSecretRef struct {
	// +optional
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	Name *string `json:"name"`

	// +optional
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	Namespace *string `json:"namespace"`
}

// +kubebuilder:object:generate=true
type Config struct {
	// +kubebuilder:validation:Required
	Env *[]EnvItem `json:"env" validate:"required"`

	// +kubebuilder:validation:Required
	Secrets *[]Secret `json:"secrets" validate:"required"`
}

// +kubebuilder:object:generate=true
type EnvItem struct {
	// +kubebuilder:validation:Required
	Name string `json:"name" validate:"required"`

	// +kubebuilder:validation:Required
	Value string `json:"value" validate:"required"`
}

// +kubebuilder:object:generate=true
type Secret struct {
	// +kubebuilder:validation:Required
	Name string `json:"name" validate:"required"`

	// +kubebuilder:validation:Required
	Type string `json:"type" validate:"required"`

	// +kubebuilder:validation:Required
	VaultPath string `json:"vaultPath" validate:"required"`

	// +optional
	DestinationSecretName *string `json:"destinationSecretName"`
}

// +kubebuilder:object:generate=true
type AdditionalCatalogSource struct {
	// Name of the additional catalog source
	// +kubebuilder:validation:Pattern=`^[a-z]([-a-z0-9]*[a-z0-9])?$`
	Name string `json:"name"`

	// Image url of the additional catalog source
	// +kubebuilder:validation:Required
	Image string `json:"image"`
}

// +kubebuilder:object:generate=true
type CredentialsRequest struct {
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	// Name of the credentials secret used to access cloud resources
	Name string `json:"name"`

	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	// Namespace where the credentials secret lives in the cluster
	Namespace string `json:"namespace"`

	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	// Service account name to use when authenticating
	ServiceAccount string `json:"service_account"`

	// +kubebuilder:validate:Pattern=`^[a-z0-9]{1,60}:[A-Za-z0-9]{1,60}$`
	// List of policy permissions needed to access cloud resources
	PolicyPermissions *[]string `json:"policy_permissions"`
}
