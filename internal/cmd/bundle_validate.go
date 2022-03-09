package cmd

import (
	"log"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/extractor"

	"github.com/spf13/cobra"
)

func init() {
	bundleCmd.AddCommand(validateBundleCmd)
}

var (
	validateBundleExamples = []string{
		"  # Validate a bundle given it's directory.",
		"  mtcli bundle validate <bundle_path>",
	}
	validateBundleCmd = &cobra.Command{
		Use:     "validate",
		Short:   "Validate a bundle given it's directory.",
		Long:    "Same as `$ opm alpha bundle validate <image>` but works locally.",
		Example: strings.Join(validateBundleExamples, "\n"),
		Args:    cobra.ExactArgs(1),
		Run:     validateBundleMain,
	}
)

func validateBundleMain(cmd *cobra.Command, args []string) {
	path := args[0]
	bundleExtractor := extractor.NewBundleExtractor()
	if err := bundleExtractor.ValidateBundle(nil, path); err != nil {
		log.Fatalf("Failed to validate bundle %s: %s.", path, err)
	}
}
