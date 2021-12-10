package v1

import (
	"encoding/json"

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

//+kubebuilder:object:generate=true
// using list so we can easily DeepCopy
type AddOnParameterList struct {
	Items []AddOnParameter `json:"items"`
}

func (a *AddOnParameterList) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &a.Items)
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
// using list so we can easily DeepCopy
type AddOnRequirementList struct {
	Items []AddOnRequirement `json:"items"`
}

func (a *AddOnRequirementList) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &a.Items)
}

//+kubebuilder:object:generate=true
type AddOnSubOperator struct {
	OperatorName      string `json:"operator_name" validate:"required"`
	OperatorNamespace string `json:"operator_namespace" validate:"required"`
	Enabled           bool   `json:"enabled" validate:"required"`
}

//+kubebuilder:object:generate=true
// using list so we can easily DeepCopy
type AddOnSubOperatorList struct {
	Items []AddOnSubOperator `json:"items"`
}

func (a AddOnSubOperatorList) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &a.Items)
}
