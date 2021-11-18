package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"

	"github.com/spf13/cobra"
)

func init() {
	mtcli.AddCommand(listBundlesCmd)
}

var (
	listBundlesExamples = []string{
		"  #List all the bundles present in an index image.",
		"  mtcli list-bundles <index_image>",
	}
	listBundlesCmd = &cobra.Command{
		Use:     "list-bundles",
		Short:   "List all the bundles present in an index image.",
		Example: strings.Join(validateExamples, "\n"),
		Args:    cobra.ExactArgs(1),
		Run:     listBundlesMain,
	}
)

func listBundlesMain(cmd *cobra.Command, args []string) {
	indexImageUrl := args[0]
	allBundles, err := utils.ExtractAndParseAllAddons(indexImageUrl)
	if err != nil {
		log.Fatalf("Failed to extract and parse bundles from the given index image: Error: %s \n", err.Error())
	}

	for _, bundle := range allBundles {
		version, err := bundle.Version()
		if err != nil {
			log.Fatalf("Failed to parse bundle: %s. Error: %s", bundle.Name, err.Error())
		}
		fmt.Printf("%s.v%s\n", bundle.Name, version)
	}
}
