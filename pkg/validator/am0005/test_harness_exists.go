package am0005

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	imageparser "github.com/novln/docker-parser"
)

const (
	code            = 5
	name            = "test_harness"
	description     = "Ensure that an addon has a valid testharness image"
	quayRegistryApi = "https://quay.io/v2"
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
	}, nil
}

type TestHarnessExists struct {
	*validator.Base
}

func (t *TestHarnessExists) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	res, err := imageparser.Parse(mb.AddonMeta.TestHarness)
	if err != nil {
		return t.Fail("Failed to parse testharness url")
	}
	if res.Registry() != "quay.io" {
		return t.Fail("Testharness image is not in the quay.io registry")
	}
	return t.checkImageExists(res)
}

func (t *TestHarnessExists) checkImageExists(imageUri *imageparser.Reference) validator.Result {
	apiUrl := fmt.Sprintf("%s/%s/manifests/%s", quayRegistryApi, imageUri.ShortName(), imageUri.Tag())
	resp, err := http.Get(apiUrl)
	if err != nil {
		if isRetryable(err) {
			return t.RetryableError(err)
		}
		return t.Error(err)
	}

	// Retry on 5xx responses
	if resp.StatusCode >= 500 {
		return t.RetryableError(errors.New("received 5XX response from quay.io"))
	}

	if resp.StatusCode == 200 {
		return t.Success()
	}
	return t.Fail(fmt.Sprintf("Test harness image doesn't exist. Received non 200 response code from quay: '%v'.", resp.StatusCode))
}

func isRetryable(err error) bool {
	urlErr, ok := err.(*url.Error)
	return ok && urlErr.Timeout()
}
