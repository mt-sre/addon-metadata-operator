package v1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
Please keep in sync with managed-tenants-cli schemas and OCM API spec:
- https://github.com/mt-sre/managed-tenants-cli/blob/main/managedtenants/data/metadata.schema.yaml
- https://github.com/mt-sre/managed-tenants-cli/blob/main/managedtenants/data/imageset.schema.yaml
- <redhat_internal_gitlab>/service/uhc-clusters-service/pkg/api/addons.go

We need to generate deepcopy methods for all complex types. Using the non-root
generation annotation as these types don't implement runtime.Object interface.


Update zz_generated.deepcopy.go with:
	$ make generate
*/

//+kubebuilder:object:generate=true
type AddOnParameter struct {
	ID           string                      `json:"id" validate:"required"`
	Name         string                      `json:"name" validate:"required"`
	Description  string                      `json:"description" validate:"required"`
	ValueType    AddOnParameterValueType     `json:"value_type" validate:"required"`
	Validation   *string                     `json:"validation"`
	Required     bool                        `json:"required" validate:"required"`
	Editable     bool                        `json:"editable" validate:"required"`
	Enabled      bool                        `json:"enabled" validate:"required"`
	DefaultValue *string                     `json:"default_value"`
	Order        *int                        `json:"order"`
	Options      *[]AddOnParameterOption     `json:"options"`
	Conditions   *[]AddOnResourceRequirement `json:"conditions"`
}

type AddOnParameterValueType string

const (
	AddOnParameterValueTypeString   AddOnParameterValueType = "string"
	AddOnParameterValueTypeBoolean  AddOnParameterValueType = "boolean"
	AddOnParameterValueTypeNumber   AddOnParameterValueType = "number"
	AddOnParameterValueTypeCIDR     AddOnParameterValueType = "cidr"
	AddOnParameterValueTypeResource AddOnParameterValueType = "resource"
)

//+kubebuilder:object:generate=true
type AddOnParameterOption struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}

//+kubebuilder:object:generate=true
type AddOnResourceRequirement struct {
	Resource AddOnRequirementResourceType    `json:"resource" validate:"required"`
	Data     AddOnRequirementData            `json:"data" validate:"required"`
	Status   *AddOnResourceRequirementStatus `json:"status"`
}

type AddOnRequirementData map[string]apiextensionsv1.JSON

//+kubebuilder:object:generate=true
type AddOnResourceRequirementStatus struct {
	Fulfilled *bool    `json:"fulfilled"`
	ErrorMsgs []string `json:"error_msgs"`
}

type AddOnRequirementResourceType string

const (
	AddOnRequirementResourceTypeCluster     AddOnRequirementResourceType = "cluster"
	AddOnRequirementResourceTypeAddOn       AddOnRequirementResourceType = "addon"
	AddOnRequirementResourceTypeMachinePool AddOnRequirementResourceType = "machine_pool"
)

//+kubebuilder:object:generate=true
type AddOnRequirement struct {
	ID       string                          `json:"id" validate:"required"`
	Resource AddOnRequirementResourceType    `json:"resource" validate:"required"`
	Data     AddOnRequirementData            `json:"data" validate:"required"`
	Status   *AddOnResourceRequirementStatus `json:"status"`
	Enabled  bool                            `json:"enabled" validate:"required"`
}

//+kubebuilder:object:generate=true
type AddOnSubOperator struct {
	OperatorName      string `json:"operator_name" validate:"required"`
	OperatorNamespace string `json:"operator_namespace" validate:"required"`
	Enabled           bool   `json:"enabled" validate:"required"`
}

// +kubebuilder:validation:Pattern=`^([A-Za-z -]+ <[0-9A-Za-z_.-]+@redhat\.com>,?)+$`
type Notification string

//+kubebuilder:object:generate=true
type Monitoring struct {
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`

	// +kubebuilder:validation:Required
	MatchNames []string `json:"matchNames"`

	// +kubebuilder:validation:Required
	MatchLabels map[string]string `json:"matchLabels"`
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
	EscalationPolicy string `json:"snitchNamePostFix"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	AcknowledgeTimeout int `json:"acknowledgeTimeout"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	ResolveTimeout int `json:"resolveTimeout"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	SecretName string `json:"secretName"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	SecretNamespace string `json:"secretNamespace"`
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
	Tags []Tag `json:"tags"`
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
	Env *[]EnvItem `json:"env"`
}

//+kubebuilder:object:generate=true
type EnvItem struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	Value string `json:"value"`
}
