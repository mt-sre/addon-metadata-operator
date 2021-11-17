package v1

import "encoding/json"

/*
Please keep in sync with managed-tenants-cli schema:
https://github.com/mt-sre/managed-tenants-cli/blob/main/managedtenants/data/metadata.schema.yaml

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
	Options      *[]AddOnParameterOption     `json:"options"`
	Conditions   *[]AddOnParameterConditions `json:"conditions"`
}

//+kubebuilder:object:generate=true
// using list so we can easily DeepCopy
type AddOnParameterList struct {
	Items []AddOnParameter `json:"items"`
}

func (a AddOnParameterList) UnmarshalJSON(b []byte) error {
	res := AddOnParameterList{}
	return json.Unmarshal(b, &res.Items)
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
type AddOnParameterConditions struct {
	Resource AddOnParameterResourceType `json:"resource" validate:"required"`
	Data     AddOnParameterData         `json:"data" validate:"required"`
}

type AddOnParameterResourceType string

const (
	AddOnParameterResourceTypeCluster AddOnParameterResourceType = "cluster"
)

//+kubebuilder:object:generate=true
type AddOnParameterData struct {
	AWSStsEnabled   *bool     `json:"aws.sts.enabled"`
	CCSEnabled      *bool     `json:"ccs.enabled"`
	CloudProviderID *[]string `json:"cloud_provider.id"`
	ProductID       *[]string `json:"product.id"`
	VersionRawID    *string   `json:"version.raw_id"`
}

//+kubebuilder:object:generate=true
type AddOnRequirement struct {
	ID       string                       `json:"id" validate:"required"`
	Resource AddOnRequirementResourceType `json:"resource" validate:"required"`
	Data     AddOnRequirementData         `json:"data" validate:"required"`
	Enabled  bool                         `json:"enabled" validate:"required"`
}

//+kubebuilder:object:generate=true
// using list so we can easily DeepCopy
type AddOnRequirementList struct {
	Items []AddOnRequirement `json:"items"`
}

func (a AddOnRequirementList) UnmarshalJSON(b []byte) error {
	res := AddOnRequirementList{}
	return json.Unmarshal(b, &res.Items)
}

type AddOnRequirementResourceType string

const (
	AddOnRequirementResourceTypeCluster     AddOnRequirementResourceType = "cluster"
	AddOnRequirementResourceTypeAddOn       AddOnRequirementResourceType = "addon"
	AddOnRequirementResourceTypeMachinePool AddOnRequirementResourceType = "machine_pool"
)

//+kubebuilder:object:generate=true
type AddOnRequirementData struct {
	ID                      *string   `json:"id"`
	State                   *string   `json:"state"`
	AWSStsEnabled           *bool     `json:"aws.sts.enabled"`
	CloudProviderID         *[]string `json:"cloud_provider.id"`
	ProductID               *[]string `json:"product.id"`
	ComputeMemory           *int      `json:"compute.memory"`
	ComputeCPU              *int      `json:"compute.cpu"`
	CCSEnabled              *bool     `json:"ccs.enabled"`
	NodesCompute            *int      `json:"nodes.compute"`
	NodesComputeMachineType *[]string `json:"nodes.compute_machine_type.id"`
	VersionRawID            *string   `json:"version.raw_id"`
	InstanceType            *[]string `json:"instance_type"`
	Replicas                *int      `json:"replicas"`
}

//+kubebuilder:object:generate=true
type AddOnSubOperator struct {
	OperatorName      string `json:"operatorName" validate:"required"`
	OperatorNamespace string `json:"operatorNamespace" validate:"required"`
	Enabled           bool   `json:"enabled" validate:"required"`
}

//+kubebuilder:object:generate=true
// using list so we can easily DeepCopy
type AddOnSubOperatorList struct {
	Items []AddOnSubOperator `json:"items"`
}

func (a AddOnSubOperatorList) UnmarshalJSON(b []byte) error {
	res := AddOnSubOperatorList{}
	return json.Unmarshal(b, &res.Items)
}
