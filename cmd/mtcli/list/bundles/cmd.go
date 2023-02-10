package bundles

import (
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/extractor"

	"github.com/spf13/cobra"
)

func examples() string {
	return strings.Join([]string{
		"  #List all the bundles present in an index image.",
		"  mtcli list bundles <index_image>",
	}, "\n")
}

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:     "bundles",
		Short:   "List all the bundles present in an index image.",
		Example: examples(),
		Args:    cobra.ExactArgs(1),
		RunE:    run,
	}
}

func run(cmd *cobra.Command, args []string) error {
	indexImageURL := args[0]

	extractor := extractor.New()
	allBundles, err := extractor.ExtractAllBundles(cmd.Context(), indexImageURL)
	if err != nil {
		return fmt.Errorf("extracting and parsing bundles from index image %q: %w", indexImageURL, err)
	}

	var operatorVersionedNames []string
	for _, bundle := range allBundles {
		csv := bundle.ClusterServiceVersion

		operatorVersionedNames = append(operatorVersionedNames, csv.Name)
	}

	fmt.Fprintln(cmd.OutOrStdout(), strings.Join(operatorVersionedNames, "\n"))

	return nil
}
