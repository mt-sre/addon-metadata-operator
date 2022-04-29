# Adding Validators

## Overview

New validation logic for addon metadata can be added by creating a validator
under the [validator](../pkg/validator) package. Please check however to see
if the validation you are implementing already fits under an existing validator.
In that case it would be appropriate to just extend the functionality of the
exising validator instead of creating a new one.

## Getting Started

First determine the unique code to assign to the validator. This code will
be selected based on the existing validators and any validators that are
still works in progress. The code should equal `n+1` where `n` is the
current greatest validator code. Then choose a meaningful name and
description to describe the validator to end users as this is what
they will see when a validation result is returned to them.

## Creating a new package

Under the [validator](../pkg/validator) package create a subpackage named
for the unique code of your validator. For example, a validator with code
`AM9999` will have a full import path of
`github.com/mt-sre/addon-metadata-operator/pkg/validator/am9999`.
Within this package you should have at least two `.go` files. One file
will implement the validator and the second will implement unit tests for
that validator.

## Implementing the required interfaces

Every validator must implement the `validator.Validator` interface:

```go
type Validator interface {
	Code() Code
	Name() string
	Description() string
	Run(context.Context, types.MetaBundle) Result
}
```

This is partially achieved by including a `*validator.Base` as a
member of your validator. However the `Run` method is up to you
to define. Luckily, the `*validator.Base` helps with that by
providing some helper functions
(`Success`, `Fail`, `Error`, `RetryableError`) that ensures you
return a proper `validator.Result` based on the logic of
your validator.

### Initializers

In addition to the validator itself your package must provide
an initializer `type Initializer func(Dependencies) (Validator, error)`
which injects dependencies from a `validator.Runner` instance
and either returns an initialized `validator.Validator` or an error.
This initializer must then be passed to the `validator.Register` function
which should be called at the top of your validator implementation file.

## Testing

Once your validator is implemented a minimum of two tests are required.
One test should validate at least one valid `types.MetaBundle` that
will result in a success. A second test should do the same, but with
invalid bundles which produce a fail result. The
[testutils](../pkg/validator/testutils/) package provides some helpers
to assist with this testing. Calling `testutils.DefaultValidBundleMap`
will provide bundles which should always pass validation. Additionally,
calling `testutils.NewValidatorTester` will initialize your validator
with configuratble dependencies and offers two methods,
`TestValidBundles` and `TestInvalidBundles` which both accept a map
of descriptive names to bundles and ensure they are valid or invalid
respectively.

## Enabling your validator

Once your validator is ready it must be imported within the
[register](../pkg/validator/register) package as a blank import.
Any packages which then import the `register` package will register
the included validators with a `validator.Runner` instance. This makes
it possible to isolate side-effectual imports during testing or other
unwanted scenarios.

## Running tests

Unit tests can be run using `./mage test:unit` and the `go test` command
can be used to run specific tests directly.
