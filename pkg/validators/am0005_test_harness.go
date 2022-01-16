package validators

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	imageparser "github.com/novln/docker-parser"
)

const quayRegistryApi = "https://quay.io/v2"

var AM0005 = types.NewValidator(
	"AM0005",
	types.ValidateFunc(ValidateTestHarness),
	types.ValidatorName("test_harness"),
	types.ValidatorDescription("Ensure that an addon has a valid testharness image"),
)

func init() {
	Registry.Add(AM0005)
}

func ValidateTestHarness(cfg types.ValidatorConfig, mb types.MetaBundle) types.ValidatorResult {
	res, err := imageparser.Parse(mb.AddonMeta.TestHarness)
	if err != nil {
		return Fail("Failed to parse testharness url")
	}
	if res.Registry() != "quay.io" {
		return Fail("Testharness image is not in the quay.io registry")
	}
	return checkImageExists(res)
}

func checkImageExists(imageUri *imageparser.Reference) types.ValidatorResult {
	apiUrl := fmt.Sprintf("%s/%s/manifests/%s", quayRegistryApi, imageUri.ShortName(), imageUri.Tag())
	resp, err := http.Get(apiUrl)
	if err != nil {
		if isRetryable(err) {
			return RetryableError(err)
		}
		return Error(err)
	}

	// Retry on 5xx responses
	if resp.StatusCode >= 500 {
		return RetryableError(errors.New("Retrying 500 status code response from quay.io."))
	}

	if resp.StatusCode == 200 {
		return Success()
	}
	return Fail(fmt.Sprintf("Test harness image doesn't exist. Received non 200 response code from quay: '%v'.", resp.StatusCode))
}

func isRetryable(err error) bool {
	urlErr, ok := err.(*url.Error)
	return ok && urlErr.Timeout()
}
