package validator

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/go-logr/logr"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
)

var initializers []Initializer

// Register queues an "Initializer" which the "Runner"
// will invoke upon startup.
func Register(init Initializer) {
	initializers = append(initializers, init)
}

// Initializer is a function which will initialize a Validator
// with dependencies. An error is returned if the Validator
// cannot be initalized properly.
type Initializer func(Dependencies) (Validator, error)

// Dependencies abstracts common dependencies for Validators.
type Dependencies struct {
	Logger     logr.Logger
	OCMClient  OCMClient
	QuayClient QuayClient
}

// NewRunner returns a Runner configured with a variadic
// slice of options or an error if an issue occurs.
func NewRunner(opts ...RunnerOption) (*Runner, error) {
	var cfg RunnerConfig

	cfg.Option(opts...)

	if err := cfg.Default(); err != nil {
		return nil, err
	}

	deps := Dependencies{
		Logger:     cfg.Logger,
		OCMClient:  cfg.OCMClient,
		QuayClient: cfg.QuayClient,
	}

	entries := make(map[Code]validatorEntry)

	for _, init := range cfg.Initializers {
		val, err := init(deps)
		if err != nil {
			return nil, err
		}

		if existing, ok := entries[val.Code()]; ok {
			return nil, fmt.Errorf(
				"code '%d' is already registered for validator '%s'",
				val.Code(),
				existing.Name(),
			)
		}

		entries[val.Code()] = validatorEntry{
			Validator: val,
		}
	}

	return &Runner{
		cfg:     cfg,
		entries: entries,
	}, nil
}

type Runner struct {
	cfg     RunnerConfig
	entries map[Code]validatorEntry
}

func (r *Runner) Run(ctx context.Context, mb types.MetaBundle, filters ...Filter) <-chan Result {
	resultCh := make(chan Result)

	var wg sync.WaitGroup

	vals := r.GetValidators(filters...)

	wg.Add(len(vals))

	for _, val := range vals {
		go func(v Validator) {
			defer wg.Done()

			run := r.applyMiddleware(v.Run)

			select {
			case <-ctx.Done():
			case resultCh <- run(ctx, mb):
			}
		}(val)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	return resultCh
}

func (r *Runner) GetValidators(filters ...Filter) []Validator {
	var result ValidatorList

	for _, e := range r.entries {
		if !e.Satisfies(filters...) {
			continue
		}

		result = append(result, e.Validator)
	}

	sort.Sort(result)

	return result
}

func (r *Runner) applyMiddleware(run RunFunc) RunFunc {
	res := run

	for _, mw := range r.cfg.Middleware {
		res = mw.Wrap(res)
	}

	return res
}

func (r *Runner) CleanUp() error {
	return r.cfg.OCMClient.CloseConnection()
}

type RunnerConfig struct {
	Initializers []Initializer
	Logger       logr.Logger
	Middleware   []Middleware
	OCMClient    OCMClient
	QuayClient   QuayClient
}

func (c *RunnerConfig) Option(opts ...RunnerOption) {
	for _, opt := range opts {
		opt.ApplyToRunnerConfig(c)
	}
}

func (c *RunnerConfig) Default() error {
	if len(c.Initializers) == 0 {
		c.Initializers = initializers
	}

	if c.Logger.GetSink() == nil {
		c.Logger = logr.Discard()
	}

	if c.OCMClient == nil {
		o, err := NewDefaultOCMClient()
		if err != nil {
			if !IsOCMClientAuthError(err) {
				return err
			}

			c.Logger.Error(err, "setting default configs:")
		}

		c.OCMClient = o
	}

	if c.QuayClient == nil {
		c.QuayClient = NewQuayClient()
	}

	return nil
}

type RunnerOption interface {
	ApplyToRunnerConfig(*RunnerConfig)
}

type WithLogger struct{ logr.Logger }

func (l WithLogger) ApplyToRunnerConfig(c *RunnerConfig) { c.Logger = l.Logger }

type WithInitializers []Initializer

func (i WithInitializers) ApplyToRunnerConfig(c *RunnerConfig) { c.Initializers = i }

type WithMiddleware []Middleware

func (m WithMiddleware) ApplyToRunnerConfig(c *RunnerConfig) { c.Middleware = m }

type WithOCMClient struct{ OCMClient }

func (o WithOCMClient) ApplyToRunnerConfig(c *RunnerConfig) { c.OCMClient = o }

type WithQuayClient struct{ QuayClient }

func (q WithQuayClient) ApplyToRunnerConfig(c *RunnerConfig) { c.QuayClient = q }

type validatorEntry struct {
	Validator
}

func (e *validatorEntry) Satisfies(filters ...Filter) bool {
	if len(filters) == 0 {
		return true
	}

	for _, f := range filters {
		if f == nil || f(e) {
			continue
		}

		return false
	}

	return true
}

type Filter func(Validator) bool

func MatchesCodes(codes ...Code) Filter {
	return func(v Validator) bool {
		for _, c := range codes {
			if v.Code() == c {
				return true
			}
		}

		return false
	}
}

func Not(f Filter) Filter {
	return func(v Validator) bool {
		return !f(v)
	}
}
