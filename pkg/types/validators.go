package types

import (
	"context"
	"time"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/operator-framework/operator-registry/pkg/registry"
)

// NewValidator returns a configured validator.
// code is a required parameter identifying the "AMXXXX" identifier for a validator
// runner is the validation process the validator will run.
// opts is a slice of additional options which can be configured.
func NewValidator(code string, runner ValidatorRunner, opts ...ValidatorOption) Validator {
	val := Validator{
		Runner: runner,
		ValidatorConfig: ValidatorConfig{
			Code: code,
		},
	}

	val.Options(opts...)

	return val
}

// Validator encapsulates display fields and implementation for
// tasks performed against MetaBundles to validate their data.
type Validator struct {
	Runner ValidatorRunner
	ValidatorConfig
}

// Validate runs a validator's Runner against a MetaBundle.
func (v *Validator) Validate(mb MetaBundle) ValidatorResult {
	return v.Runner.Run(v.ValidatorConfig, mb)
}

// Option applies an option to an instance of Validator.
func (v *Validator) Options(opts ...ValidatorOption) {
	for _, opt := range opts {
		opt(&v.ValidatorConfig)
	}
}

// ApplyMiddleware wraps a Validator's Runner with the supplied slice of Middleware.
// Note that the order of the Middleware is important.
func (v *Validator) ApplyMiddleware(mws ...Middleware) {
	for _, mw := range mws {
		v.Runner = mw.Wrap(v.Runner)
	}
}

// ValidatorConfig encapsulates configruation data for a validator.
type ValidatorConfig struct {
	Code        string
	Name        string
	Desc        string
	Ctx         context.Context
	RetryParams *RetryParams
}

// RetryParms encapsulates parameters used by a validator to conifgure its retries.
type RetryParams struct {
	Count   int
	Delayer Delayer
}

// ValidatorOption applis options to an instance of ValidatorConfig.
type ValidatorOption func(c *ValidatorConfig)

// ValidatorName sets the display name of the validator.
func ValidatorName(name string) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.Name = name
	}
}

// ValidatorDescription sets the display description of the validator
func ValidatorDescription(desc string) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.Desc = desc
	}
}

// ValidatorContext sets the Ctx of an instance of ValidatorConfig.
func ValidatorContext(ctx context.Context) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.Ctx = ctx
	}
}

// ValidatorRetryParams sets the RetryParams of an instance of ValidatorConfig.
func ValidatorRetryParams(params *RetryParams) ValidatorOption {
	return func(c *ValidatorConfig) {
		c.RetryParams = params
	}
}

// ValidatorRunner describes objects which run validations against
// a MetaBundle and produce a ValidatorResult. This is analogous to
// the http.Handler interface.
type ValidatorRunner interface {
	// Run validates a MetaBundle and produces a ValidatorResult
	// indicating whether the validations succedded, failed, or errored out.
	// This is analogous to Handler.ServeHTTP
	Run(ValidatorConfig, MetaBundle) ValidatorResult
}

// ValidateFunc wraps functions which validate MetaBundles.
type ValidateFunc func(ValidatorConfig, MetaBundle) ValidatorResult

func (f ValidateFunc) Run(cfg ValidatorConfig, mb MetaBundle) ValidatorResult {
	return f(cfg, mb)
}

// Delayer abstracts behavior used to generate successive delay values.
type Delayer interface {
	// Delay returns a time.Duration value algorithmically based
	// on the implementation of a particular Delayer.
	Delay() time.Duration
}

// Middleware abstracts behavior which extends types.ValidatorRunner
// functionality by executing additional steps before/after validations
// are run such as adding retries or logging.
type Middleware interface {
	// Wrap applies the implemented middleware to an instance of types.ValidatorRunner
	Wrap(ValidatorRunner) ValidatorRunner
}

// ValidatorList - implements Sort interface to sort validators per Code
type ValidatorList []Validator

func (v ValidatorList) Len() int {
	return len(v)
}

func (v ValidatorList) Less(i, j int) bool {
	return v[i].Code < v[j].Code
}

func (v ValidatorList) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

type ValidatorTest interface {
	Name() string
	Run(MetaBundle) ValidatorResult
	SucceedingCandidates() ([]MetaBundle, error)
	FailingCandidates() ([]MetaBundle, error)
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

func (vr ValidatorResult) IsRetryableError() bool {
	return vr.Error != nil && vr.RetryableError
}
