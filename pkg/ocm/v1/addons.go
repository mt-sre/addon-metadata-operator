package v1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
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
	// +kubebuilder:validation:Required
	ID string `json:"id" validate:"required"`

	// +kubebuilder:validation:Required
	Name string `json:"name" validate:"required"`

	// +kubebuilder:validation:Required
	Description string `json:"description" validate:"required"`

	// +kubebuilder:validation:Required
	ValueType AddOnParameterValueType `json:"value_type" validate:"required"`

	// +optional
	Validation *string `json:"validation"`
	// +kubebuilder:validation:Required
	Required bool `json:"required" validate:"required"`

	// +kubebuilder:validation:Required
	Editable bool `json:"editable" validate:"required"`

	// +kubebuilder:validation:Required
	Enabled bool `json:"enabled" validate:"required"`

	// +optional
	DefaultValue *string `json:"default_value"`

	// +optional
	Order *int `json:"order"`

	// +optional
	Options *[]AddOnParameterOption `json:"options"`

	// +optional
	Conditions *[]AddOnResourceRequirement `json:"conditions"`
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
	// +kubebuilder:validation:Required
	Name string `json:"name" validate:"required"`

	// +kubebuilder:validation:Required
	Value string `json:"value" validate:"required"`
}

//+kubebuilder:object:generate=true
type AddOnResourceRequirement struct {
	// +kubebuilder:validation:Required
	Resource AddOnRequirementResourceType `json:"resource" validate:"required"`

	// +kubebuilder:validation:Required
	Data AddOnRequirementData `json:"data" validate:"required"`

	// +optional
	Status *AddOnResourceRequirementStatus `json:"status"`
}

type AddOnRequirementData map[string]apiextensionsv1.JSON

//+kubebuilder:object:generate=true
type AddOnResourceRequirementStatus struct {
	// +optional
	Fulfilled *bool `json:"fulfilled"`

	// +optional
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
	// +kubebuilder:validation:Required
	ID string `json:"id" validate:"required"`

	// +kubebuilder:validation:Required
	Resource AddOnRequirementResourceType `json:"resource" validate:"required"`

	// +kubebuilder:validation:Required
	Data AddOnRequirementData `json:"data" validate:"required"`

	// +optional
	Status *AddOnResourceRequirementStatus `json:"status"`

	// +kubebuilder:validation:Required
	Enabled bool `json:"enabled" validate:"required"`
}

//+kubebuilder:object:generate=true
type AddOnSubOperator struct {
	// +kubebuilder:validation:Required
	OperatorName string `json:"operator_name" validate:"required"`

	// +kubebuilder:validation:Required
	OperatorNamespace string `json:"operator_namespace" validate:"required"`

	// +kubebuilder:validation:Required
	Enabled bool `json:"enabled" validate:"required"`
}
