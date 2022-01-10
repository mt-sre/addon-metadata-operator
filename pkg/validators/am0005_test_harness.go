package validators

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	imageparser "github.com/novln/docker-parser"
)

const quayRegistryApi = "https://quay.io/v2"

var AM0005 = types.Validator{
	Code:        "AM0005",
	Name:        "test_harness",
	Description: "Ensure that an addon has a valid testharness image",
	Runner:      ValidateTestHarness,
}

func init() {
	Registry.Add(AM0005)
}

func ValidateTestHarness(mb types.MetaBundle) types.ValidatorResult {
	res, err := imageparser.Parse(mb.AddonMeta.TestHarness)
	if err != nil {
		return Fail("Failed to parse testharness url")
	}
	if res.Registry() != "quay.io" {
		return Fail("Testharness image is not in the quay.io registry")
	}
	retryCount := 5
	return checkImageExists(res, retryCount)
}

func checkImageExists(imageUri *imageparser.Reference, retryCount int) types.ValidatorResult {
	apiUrl := fmt.Sprintf("%s/%s/manifests/%s", quayRegistryApi, imageUri.ShortName(), imageUri.Tag())
	resp, err := http.Get(apiUrl)
	if err != nil {
		if retryCount > 0 && retryableError(err) {
			time.Sleep(2 * time.Second)
			return checkImageExists(imageUri, retryCount-1)
		}
		return Error(err)
	}

	if resp.StatusCode == 200 {
		return Success()
	}
	// Retry on 5xx responses
	if retryCount > 0 && resp.StatusCode >= 500 {
		time.Sleep(2 * time.Second)
		return checkImageExists(imageUri, retryCount-1)
	}
	return Fail(fmt.Sprintf("Test harness image doesn't exist. Received non 200 response code from quay: '%v'.", resp.StatusCode))
}

func retryableError(err error) bool {
	urlErr, ok := err.(*url.Error)
	return ok && urlErr.Timeout()
}
