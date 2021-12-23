package validators

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	imageparser "github.com/novln/docker-parser"
)

const quayRegistryApi = "https://quay.io/v2"

var AM0005 = utils.Validator{
	Code:        "AM0005",
	Name:        "test_harness",
	Description: "Ensure that an addon has a valid testharness image",
	Runner:      ValidateTestHarness,
}

func init() {
	Registry.Add(AM0005)
}

func ValidateTestHarness(metabundle utils.MetaBundle) (bool, string, error) {
	res, err := imageparser.Parse(metabundle.AddonMeta.TestHarness)
	if err != nil {
		return false, "Failed to parse testharness url", err
	}
	if res.Registry() != "quay.io" {
		return false, "Testharness image is not in the quay.io registry", nil
	}
	retryCount := 5
	imagePresent, err := checkImageExists(res, retryCount)
	if err != nil {
		return false, "Encountered an error when trying to access the remote image", err
	}
	if !imagePresent {
		return false, "Test harness image doesnt exist", nil
	}
	return true, "", nil
}

func checkImageExists(imageUri *imageparser.Reference, retryCount int) (bool, error) {
	apiUrl := fmt.Sprintf("%s/%s/manifests/%s", quayRegistryApi, imageUri.ShortName(), imageUri.Tag())
	resp, err := http.Get(apiUrl)
	if err != nil {
		if retryCount > 0 && retryableError(err) {
			time.Sleep(2 * time.Second)
			return checkImageExists(imageUri, retryCount-1)
		}
		return false, err
	}

	if resp.StatusCode == 200 {
		return true, nil
	}
	// Retry on 5xx responses
	if retryCount > 0 && resp.StatusCode >= 500 {
		time.Sleep(2 * time.Second)
		return checkImageExists(imageUri, retryCount-1)
	}
	return false, nil
}

func retryableError(err error) bool {
	urlErr, ok := err.(*url.Error)
	return ok && urlErr.Timeout()
}
