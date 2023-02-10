package validate

import (
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/extractor"

	"github.com/spf13/cobra"
)

func examples() string {
	return strings.Join([]string{
		"  # Validate a bundle given it's directory.",
		"  mtcli bundle validate <bundle_path>",
	}, "\n")
}

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:     "validate",
		Short:   "Validate a bundle given it's directory.",
		Long:    "Same as `$ opm alpha bundle validate <image>` but works locally.",
		Example: examples(),
		Args:    cobra.ExactArgs(1),
		RunE:    run,
	}
}

func run(cmd *cobra.Command, args []string) error {
	path := args[0]
	bundleExtractor := extractor.NewBundleExtractor()

	if err := bundleExtractor.ValidateBundle(cmd.Context(), nil, path); err != nil {
		return fmt.Errorf("validating bundle %s: %w", path, err)
	}

	return nil
}
