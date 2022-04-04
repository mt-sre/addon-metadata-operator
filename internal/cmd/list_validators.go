package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	_ "github.com/mt-sre/addon-metadata-operator/pkg/validator/register"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.AddCommand(listValidatorsCmd)
}

var (
	listValidatorsExamples = []string{
		"  # List all the registered validators.",
		"  mtcli list validators",
	}
	listValidatorsCmd = &cobra.Command{
		Use:     "validators",
		Short:   "List all the registered validators.",
		Example: strings.Join(listValidatorsExamples, "\n"),
		Run:     listValidatorsMain,
	}
)

func listValidatorsMain(cmd *cobra.Command, args []string) {
	t := simpletable.New()
	t.Header = getTableHeaders()

	runner, err := validator.NewRunner()
	if err != nil {
		fmt.Printf("Unable to list validator: %s\n", err)
		os.Exit(1)
	}

	for _, v := range runner.GetValidators() {
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

func newTableRow(v validator.Validator) []*simpletable.Cell {
	return []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: v.Code().String()},
		{Align: simpletable.AlignLeft, Text: v.Name()},
		{Align: simpletable.AlignLeft, Text: v.Description()},
	}
}
