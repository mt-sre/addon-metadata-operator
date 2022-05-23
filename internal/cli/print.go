package cli

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/fatih/color"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

var (
	green            = color.New(color.FgGreen).SprintFunc()
	red              = color.New(color.FgRed).SprintFunc()
	intenselyBoldRed = color.New(color.Bold, color.FgHiRed).SprintFunc()
)

func NewResultTable() ResultTable {
	var table ResultTable

	table.Table = simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "STATUS"},
			{Align: simpletable.AlignCenter, Text: "CODE"},
			{Align: simpletable.AlignCenter, Text: "NAME"},
			{Align: simpletable.AlignCenter, Text: "DESCRIPTION"},
			{Align: simpletable.AlignCenter, Text: "FAILURE MESSAGE"},
		},
	}

	table.SetStyle(simpletable.StyleCompactLite)

	return table
}

type ResultTable struct {
	*simpletable.Table
}

func (t *ResultTable) WriteRow(row []*simpletable.Cell) {
	t.Body.Cells = append(t.Body.Cells, row)
}

func (t *ResultTable) WriteResult(res validator.Result) {
	row := resultToRow(res)

	if res.IsSuccess() {
		t.WriteRow(append(row, &simpletable.Cell{Align: simpletable.AlignLeft, Text: "None"}))
	} else if res.IsError() {
		t.WriteRow(append(row, &simpletable.Cell{Align: simpletable.AlignLeft, Text: res.Error.Error()}))
	} else {
		for _, msg := range res.FailureMsgs {
			t.WriteRow(append(row, &simpletable.Cell{Align: simpletable.AlignLeft, Text: msg}))
		}
	}
}

func resultToRow(res validator.Result) []*simpletable.Cell {
	var status string

	if res.IsSuccess() {
		status = green("Success")
	} else if res.IsError() {
		status = intenselyBoldRed("Error")
	} else {
		status = red("Failed")
	}

	return []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: status},
		{Align: simpletable.AlignLeft, Text: res.Code.String()},
		{Align: simpletable.AlignLeft, Text: res.Name},
		{Align: simpletable.AlignLeft, Text: res.Description},
	}
}

// PrintValidationErrors - helper to pretty print validationErrors
func PrintValidationErrors(errs []error) {
	fmt.Printf("\n%s\n", red("Failed with the following errors:"))
	for _, err := range errs {
		fmt.Printf("%s\n", err.Error())
	}
}
