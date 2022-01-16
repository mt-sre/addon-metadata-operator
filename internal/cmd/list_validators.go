package cmd

import (
	"fmt"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
	"github.com/spf13/cobra"
)

func init() {
	mtcli.AddCommand(listValidatorsCmd)
}

var (
	listValidatorsExamples = []string{
		"  # List all the registered validators.",
		"  mtcli list-validators",
	}
	listValidatorsCmd = &cobra.Command{
		Use:     "list-validators",
		Short:   "List all the registered validators.",
		Example: strings.Join(listValidatorsExamples, "\n"),
		Run:     listValidatorsMain,
	}
)

func listValidatorsMain(cmd *cobra.Command, args []string) {
	t := simpletable.New()
	t.Header = getTableHeaders()
	for _, v := range validators.Registry.ListSorted() {
		row := newTableRow(v)
		t.Body.Cells = append(t.Body.Cells, row)
	}
	t.SetStyle(simpletable.StyleCompactLite)
	fmt.Printf("%v\n\n", t.String())
}

func getTableHeaders() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "CODE"},
			{Align: simpletable.AlignCenter, Text: "NAME"},
			{Align: simpletable.AlignCenter, Text: "DESCRIPTION"},
		},
	}
}

func newTableRow(v types.Validator) []*simpletable.Cell {
	return []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: v.Code},
		{Align: simpletable.AlignLeft, Text: v.Name},
		{Align: simpletable.AlignLeft, Text: v.Desc},
	}
}
