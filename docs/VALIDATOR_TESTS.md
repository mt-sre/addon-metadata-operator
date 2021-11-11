# Validator Tests Setup Guide

This doc describes the steps to simply configure unit tests for a validator getting onboarded to this project.

## Introduction

Whenever you would add a new validator to this project, you'd also be expected to configure unit tests for that newly added validator. Don't worry, you won't be expected to code the unit tests from scratch. This project already has a generic bootstrapped setup for validator testing suite and all you would have to do is to follow some simple steps to get unit tests of your newly added validator integrated with the validators testing suite of this project.

## Background

Every validator test is represented by a dedicated type (struct) which implements the `Validator` interface (`pkg/validate/validate_test.go`)

```go
type Validator interface {
	Name() string
	Run(utils.MetaBundle) (bool, error)
	SucceedingCandidates() []utils.MetaBundle
	FailingCandidates() []utils.MetaBundle
}
```

For implementing that interface, the Validator test's struct must have four methods:
- `Name() string` : Returns the name of the validator to which the test corresponds.
- `SucceedingCandidates() []utils.MetaBundle`: Returns a slice of `MetaBundle` where each element can successfully pass the validator.
- `FailingCandidates() []utils.MetaBundle`: Returns a slice of `MetaBundle` where each element fails the validator due to any reason.
- `Run(utils.MetaBundle) (bool, error)`:  Represents running the validator over a `MetaBundle` and returning a `bool` (success: true/false) and `error` (error while executing the validator) accordingly. Succeeding Candidate `MetaBundle` fed to this method will be expected to return `(true, nil)` and Failing Candidate `MetaBundle` fed to this method will be expected to return `(false, nil/non-nil)`.

**Checkout an [Example](#example)**


## How to add a new validator test

Say, you created a new validator called `validator_default_channel` under the file `pkg/validators/validator_default_channel.go` accordingly.

* Now, go ahead and create a new file `pkg/validators/test_validator_default_channel.go`.
* Inside this file, define a struct corresponding to this validator. It can be an empty struct or non-empty depending on the way you want to implement the testing of it. For example: `type ValidateDefaultChannelTestBundle struct {}`
* Now, make this struct implement the `Validator` interface defined under `pkg/validate/validate_test.go`. (Refer to the previous section for the explanation for these methods).
* After this, your tests are ready but you have to register them to the validators' test suite now.
* For that, proceed to `pkg/validate/validate_test.go` and look for `var validatorsToTest`.
* Under it, just append the struct corresponding to the validator-test you defined under `pkg/validators/test_validator_default_channel.go`. For example, in our case, it would look like this:
```go
...
// register the validators to test here by appending them to the following slice
var validatorsToTest []Validator = []Validator{
    ValidatorAddonLabelTestBundle{},  // ignore this, it was previously present
    ValidatorDefaultChannelTestBundle{}, // pay attention! we just added this
}
...
```
* That's it! Rest will be taken care of assuming that the methods you defined for the struct are aptly coded.
* Go ahead and run the all the tests (including the validator tests) by entering `make test` and it will run:

Passing Tests
```sh
╰─ make test

[DEBUG] Go version used is go version go1.17.1 darwin/amd64
=== RUN   Test_AllValidators
=== RUN   Test_AllValidators/Addon_Label_Validator
=== PAUSE Test_AllValidators/Addon_Label_Validator
=== RUN   Test_AllValidators/Addon_Default_Channel_Validator
=== PAUSE Test_AllValidators/Addon_Default_Channel_Validator
=== CONT  Test_AllValidators/Addon_Label_Validator
=== CONT  Test_AllValidators/Addon_Default_Channel_Validator
--- PASS: Test_AllValidators (0.00s)
    --- PASS: Test_AllValidators/Addon_Label_Validator (0.00s)
    --- PASS: Test_AllValidators/Addon_Default_Channel_Validator (0.00s)
PASS
ok      github.com/mt-sre/addon-metadata-operator/internal/testutil     0.641s
```
Failing tests
```sh
╰─ make test
[DEBUG] Go version used is go version go1.17.1 darwin/amd64
ok      github.com/mt-sre/addon-metadata-operator/api/v1alpha1  0.511s
?       github.com/mt-sre/addon-metadata-operator/cmd/addon-metadata-operator   [no test files]
?       github.com/mt-sre/addon-metadata-operator/cmd/mtcli     [no test files]
?       github.com/mt-sre/addon-metadata-operator/controllers   [no test files]
?       github.com/mt-sre/addon-metadata-operator/internal/cmd  [no test files]
?       github.com/mt-sre/addon-metadata-operator/internal/testutils    [no test files]
?       github.com/mt-sre/addon-metadata-operator/pkg/utils     [no test files]
--- FAIL: Test_AllValidators (0.00s)
    --- FAIL: Test_AllValidators/Addon_Label_Validator (0.00s)
        validate_test.go:42:
                Error Trace:    validate_test.go:42
                Error:          Should be true
                Test:           Test_AllValidators/Addon_Label_Validator
        validate_test.go:42:
                Error Trace:    validate_test.go:42
                Error:          Should be true
                Test:           Test_AllValidators/Addon_Label_Validator
    --- FAIL: Test_AllValidators/Addon_Default_Channel_Validator (0.00s)
        validate_test.go:42:
                Error Trace:    validate_test.go:42
                Error:          Should be true
                Test:           Test_AllValidators/Addon_Default_Channel_Validator
        validate_test.go:42:
                Error Trace:    validate_test.go:42
                Error:          Should be true
                Test:           Test_AllValidators/Addon_Default_Channel_Validator
FAIL
FAIL    github.com/mt-sre/addon-metadata-operator/pkg/validate  0.748s
?       github.com/mt-sre/addon-metadata-operator/pkg/validators        [no test files]
FAIL
make: *** [test] Error 1
```

**Checkout an [Example](#example)**

## Example

Say, you created a validator under the file `pkg/validators/validator_default_channel.go`.

* Create `pkg/validators/test_validator_default_channel.go` containing the struct `ValidatorDefaultChannelTestBundle` corresponding to this validator-test and the methods making it implement the `Validator` interface.

<details>
  <summary><b>Click to expand:</b> `pkg/validators/test_validator_default_channel.go`</summary>

  ```go
	package validators

	import (
		"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
		"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	)

	type ValidatorDefaultChannelTestBundle struct{}

	func (val ValidatorDefaultChannelTestBundle) Name() string {
		return "Addon Default Channel Validator"
	}

	func (val ValidatorDefaultChannelTestBundle) Run(mb utils.MetaBundle) (bool, error) {
		return ValidateDefaultChannel(&mb)
	}

	func (val ValidatorDefaultChannelTestBundle) SucceedingCandidates() []utils.MetaBundle {
		return []utils.MetaBundle{
			{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					ID:             "random-operator",
					DefaultChannel: "alpha",
					Channels: []v1alpha1.Channel{
						{
							Name: "alpha",
						},
						{
							Name: "sigma",
						},
					},
				},
			},
			{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					ID:             "random-operator",
					DefaultChannel: "beta",
					Channels: []v1alpha1.Channel{
						{
							Name: "alpha",
						},
						{
							Name: "beta",
						},
					},
				},
			},
		}
	}

	func (val ValidatorDefaultChannelTestBundle) FailingCandidates() []utils.MetaBundle {
		return []utils.MetaBundle{
			{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					ID:             "random-operator",
					DefaultChannel: "alpha",
					Channels: []v1alpha1.Channel{
						{
							Name: "beta",
						},
						{
							Name: "sigma",
						},
					},
				},
			},
			{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					ID:             "random-operator",
					DefaultChannel: "beta",
					Channels: []v1alpha1.Channel{
						{
							Name: "alpha",
						},
					},
				},
			},
		}
	}
  ```
</details></br>

* Go to `pkg/validate/validate_test.go` and append the `validators.ValidatorDefaultChannelTestBundle{}` under the slice `validatorsToTest`

```go
...
// register the validators to test here by appending them to the following slice
var validatorsToTest []Validator = []Validator{
    ValidatorAddonLabelTestBundle{},  // ignore this, it was previously present
    ValidatorDefaultChannelTestBundle{}, // pay attention! we just appended this
}
...
```

* That's it! Run `make test` to see it in action :)

```sh
╰─ make test

[DEBUG] Go version used is go version go1.17.1 darwin/amd64
=== RUN   Test_AllValidators
=== RUN   Test_AllValidators/Addon_Label_Validator
=== PAUSE Test_AllValidators/Addon_Label_Validator
=== RUN   Test_AllValidators/Addon_Default_Channel_Validator
=== PAUSE Test_AllValidators/Addon_Default_Channel_Validator
=== CONT  Test_AllValidators/Addon_Label_Validator
=== CONT  Test_AllValidators/Addon_Default_Channel_Validator
--- PASS: Test_AllValidators (0.00s)
    --- PASS: Test_AllValidators/Addon_Label_Validator (0.00s)
    --- PASS: Test_AllValidators/Addon_Default_Channel_Validator (0.00s)
PASS
ok      github.com/mt-sre/addon-metadata-operator/pkg/validate
```
