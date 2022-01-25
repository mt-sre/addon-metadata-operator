package cmd

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/extractor"

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
		Example: strings.Join(listBundlesExamples, "\n"),
		Args:    cobra.ExactArgs(1),
		Run:     listBundlesMain,
	}
)

func listBundlesMain(cmd *cobra.Command, args []string) {
	indexImage := args[0]

	extractor := extractor.New()
	allBundles, err := extractor.ExtractAllBundles(indexImage)
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

	sort.Strings(operatorVersionedNames)
	fmt.Println(strings.Join(operatorVersionedNames, "\n"))
}
