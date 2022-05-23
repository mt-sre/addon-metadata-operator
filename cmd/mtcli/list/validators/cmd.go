package validators

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	_ "github.com/mt-sre/addon-metadata-operator/pkg/validator/register"
	"github.com/spf13/cobra"
)

func examples() string {
	return strings.Join([]string{
		"  # List all the registered validators.",
		"  mtcli list validators",
	}, "\n")
}

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:     "validators",
		Short:   "List all the registered validators.",
		Example: examples(),
		RunE:    run,
	}
}

func run(cmd *cobra.Command, args []string) error {
	runner, err := validator.NewRunner()
	if err != nil {
		return fmt.Errorf("listing validators: %s\n", err)
	}

	t := simpletable.New()
	t.Header = getTableHeaders()
	t.SetStyle(simpletable.StyleCompactLite)

	for _, v := range runner.GetValidators() {
		row := newTableRow(v)
		t.Body.Cells = append(t.Body.Cells, row)
	}

	fmt.Fprintf(os.Stdout, "%v\n\n", t.String())

	return nil
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
