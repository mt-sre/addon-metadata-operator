package am0005

import (
	"context"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	imageparser "github.com/novln/docker-parser"
)

const (
	code        = 5
	name        = "test_harness"
	description = "Ensure that an addon has a valid testharness image"
)

func init() {
	validator.Register(NewTestHarnessExists)
}

func NewTestHarnessExists(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(description),
	)
	if err != nil {
		return nil, err
	}

	return &TestHarnessExists{
		Base: base,
		quay: deps.QuayClient,
	}, nil
}

type TestHarnessExists struct {
	*validator.Base
	quay validator.QuayClient
}

func (t *TestHarnessExists) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	ref, err := imageparser.Parse(mb.AddonMeta.TestHarness)
	if err != nil {
		return t.Fail("Failed to parse testharness url")
	}

	if ref.Registry() != "quay.io" {
		return t.Fail("Testharness image is not in the quay.io registry")
	}

	ok, err := t.quay.HasReference(ctx, ref)
	if err != nil {
		return t.Error(err)
	}

	if !ok {
		return t.Fail(fmt.Sprintf("The testharness image %q does not exist", ref.Name()))
	}

	return t.Success()
}
