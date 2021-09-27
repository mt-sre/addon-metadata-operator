package cmd

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	addonsv1alpha1 "github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func init() {
	mtcli.AddCommand(validateCmd)
}

var (
	validateExamples = []string{
		"  # Validate an addon.yaml file on local filesystem.",
		"  mtcli validate <path/to/addon.yaml>",
		"  # Validate an addon.yaml file loaded form an URL.",
		"  mtcli validate https://<url/to/addon.yaml>",
	}
	validateCmd = &cobra.Command{
		Use:     "validate",
		Short:   "Validate an addon metadata against custom validators and the managed-tenants-cli JSON schema: https://github.com/mt-sre/managed-tenants-cli/blob/main/docs/tenants/zz_schema_generated.md.",
		Example: strings.Join(validateExamples, "\n"),
		Args:    cobra.ExactArgs(1),
		Run:     validateMain,
	}
)

func validateMain(cmd *cobra.Command, args []string) {
	addonURI := args[0]
	addonMetadata := &addonsv1alpha1.AddonMetadata{}

	data, err := readAddonMetadata(addonURI)
	log.Debugf("Raw data read from addonURI %v: \n%v\n", addonURI, string(data))

	if err != nil {
		log.Fatalf("Could not read addon metadata from URI %v, got %v.\n", addonURI, err)
	}

	if err := addonMetadata.FromYAML(data); err != nil {
		log.Fatalf("Could not load addon metadata from file %v, got %v.\n", addonURI, err)
	}

	if err := addonMetadata.Validate(); err != nil {
		addonsv1alpha1.PrintValidationErrors(err)
		log.Fatalln("Addon Metadata validation failed.")
	}
}

func readAddonMetadata(addonURI string) ([]byte, error) {
	if isLocalPath(addonURI) {
		return ioutil.ReadFile(addonURI)
	}
	if isValidURL(addonURI) {
		response, err := http.Get(addonURI)
		if err != nil {
			log.Fatalf("Could not read from URL %v, got %v.\n", addonURI, err)
		}
		defer response.Body.Close()
		return ioutil.ReadAll(response.Body)
	}
	return nil, errors.New("Invalid addon metadata URI provided.")
}

func isLocalPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isValidURL(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	return err == nil
}
