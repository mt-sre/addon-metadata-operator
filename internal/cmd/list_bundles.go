package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"

	"github.com/spf13/cobra"
)

func init() {
	listCmd.AddCommand(listBundlesCmd)
}

var (
	listBundlesExamples = []string{
		"  #List all the bundles present in an index image.",
		"  mtcli list bundles <index_image>",
	}
	listBundlesCmd = &cobra.Command{
		Use:     "bundles",
		Short:   "List all the bundles present in an index image.",
		Example: strings.Join(listBundlesExamples, "\n"),
		Args:    cobra.ExactArgs(1),
		Run:     listBundlesMain,
	}
)

func listBundlesMain(cmd *cobra.Command, args []string) {
	indexImageUrl := args[0]
	allBundles, err := utils.ExtractAndParseAddons(indexImageUrl, utils.AllAddonsIdentifier)
	if err != nil {
		log.Fatalf("Failed to extract and parse bundles from the given index image: Error: %s \n", err.Error())
	}
	var operatorVersionedNames []string
	for _, bundle := range allBundles {
		csv, err := bundle.ClusterServiceVersion()
		if err != nil {
			log.Fatalf("Failed to extract version info for bundle: %s. Error: %s", bundle.Name, err.Error())
		}
		operatorVersionedNames = append(operatorVersionedNames, csv.GetName())
	}
	fmt.Println(strings.Join(operatorVersionedNames, "\n"))
}
