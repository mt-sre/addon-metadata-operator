package validator

import (
	"context"
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

// Validator is a task which given a types.MetaBundle will
// perform checks and return a result.
type Validator interface {
	// Code returns the unique id of a Validator instance.
	Code() Code
	// Name returns the display name of a Validator instance.
	Name() string
	// Description returns the displayed description of a Validator instance.
	Description() string
	// Run executes validation tasks against a types.MetaBundle and returns the
	// result of that task. A context.Context instance is also passed to allow
	// for cancellation and timeouts to propogate through the validation task
	// and preempt any long-running processing step.
	Run(context.Context, types.MetaBundle) Result
}

// NewBase returns a base Validator implementation with a given code and optional
// parameters. An error is returned if an invalid code is given.
func NewBase(code Code, opts ...BaseOption) (*Base, error) {
	if code < Code(0) {
		return nil, fmt.Errorf("validator codes must be non-negative integers not %d", code)
	}

	cfg := Base{code: code}

	cfg.Option(opts...)
	cfg.Default()

	return &cfg, nil
}

// Base implements the base functionality used by Validator instances.
type Base struct {
	code Code
	name string
	desc string
}

func (b *Base) Code() Code          { return b.code }
func (b *Base) Name() string        { return b.name }
func (b *Base) Description() string { return b.desc }

// Option applies a variadic slice of options to a Base instance.
func (b *Base) Option(opts ...BaseOption) {
	for _, opt := range opts {
		opt(b)
	}
}

// Default applies default values for any unconfigured options.
func (b *Base) Default() {
	if b.name == "" {
		b.name = fmt.Sprintf("unnamed validator <%s>", b.code)
	}

	if b.desc == "" {
		b.desc = "no description available"
	}
}

// Success is a helper which returns a populated Success result.
func (b *Base) Success() Result {
	res := b.populateResult()
	res.success = true

	return res
}

// Fail is a helper which returns a populated Fail result.
// A variadic slice of messages are passed to describe the reason(s)
// that a validation task failed.
func (b *Base) Fail(msgs ...string) Result {
	res := b.populateResult()
	res.FailureMsgs = msgs

	return res
}

// Error is a helper which returns a populated Error result.
// An error instnace is passed to give context for what error
// caused a validation task to exit.
func (b *Base) Error(err error) Result {
	res := b.populateResult()
	res.Error = err

	return res
}

// RetryableError is a helper which returns a populated RetryableError result.
// A RetryableError indicates to middleware that the error is temporary and
// may be retried.
func (b *Base) RetryableError(err error) Result {
	res := b.populateResult()
	res.Error = err
	res.retryable = true

	return res
}

func (b *Base) populateResult() Result {
	return Result{
		Code:        b.code,
		Name:        b.name,
		Description: b.desc,
	}
}

// BaseOption abstracts functions which apply optional
// parameters to a Base instance.
type BaseOption func(*Base)

// BaseName applies the given name to a base instance.
func BaseName(name string) BaseOption {
	return func(b *Base) { b.name = name }
}

// BaseDesc applies the given description to a base instance.
func BaseDesc(desc string) BaseOption {
	return func(b *Base) { b.desc = desc }
}

// ValidatorList is a sortable slice of Validators.
type ValidatorList []Validator

func (l ValidatorList) Len() int           { return len(l) }
func (l ValidatorList) Less(i, j int) bool { return l[i].Code() < l[j].Code() }
func (l ValidatorList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

const codePrefix = "AM"

// Code is a prefixed integer ID used to distinguish Validator implementations.
type Code int

func (c Code) String() string {
	return fmt.Sprintf("%s%04d", codePrefix, c)
}

// ParseCode converts a given string to a Code value.
// An error is returned if the string is incorrectly formatted.
func ParseCode(maybeCode string) (Code, error) {
	var result Code

	if len(maybeCode) != 6 {
		return result, fmt.Errorf("code must be of the format '%sXXXX'", codePrefix)
	}

	n, err := fmt.Sscanf(strings.ToUpper(maybeCode), codePrefix+"%04d", &result)
	if err != nil || n < 1 {
		return result, fmt.Errorf("unable to parse code from '%s'", maybeCode)
	}

	return result, nil
}
