package types

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

type ValidateFunc func(mb MetaBundle) ValidatorResult

type Validator struct {
	Name        string
	Code        string
	Description string
	Runner      ValidateFunc
}

type ValidatorTest interface {
	Name() string
	Run(MetaBundle) ValidatorResult
	SucceedingCandidates() []MetaBundle
	FailingCandidates() []MetaBundle
}

type MetaBundle struct {
	AddonMeta *v1alpha1.AddonMetadataSpec
	Bundles   []registry.Bundle
}

func NewMetaBundle(addonMeta *v1alpha1.AddonMetadataSpec, bundles []registry.Bundle) *MetaBundle {
	return &MetaBundle{
		AddonMeta: addonMeta,
		Bundles:   bundles,
	}
}

// ValidatorResult - encompasses validator result information
type ValidatorResult struct {
	// true if MetaBundle validation was successful
	Success bool
	// "" if validation is successful, else information about why it failed
	FailureMsg string
	// retports error that happened in the validation code
	Error error
	// if an error occured in the validation code, determines if it was retryable
	RetryableError bool
}

func (vr ValidatorResult) IsSuccess() bool {
	return vr.Error == nil && vr.FailureMsg == ""
}

func (vr ValidatorResult) IsError() bool {
	return vr.Error != nil
}
