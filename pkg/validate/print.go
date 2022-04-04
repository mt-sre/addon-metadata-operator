package validate

import (
	"github.com/alexeyco/simpletable"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
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
		status = utils.Green("Success")
	} else if res.IsError() {
		status = utils.IntenselyBoldRed("Error")
	} else {
		status = utils.Red("Failed")
	}

	return []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: status},
		{Align: simpletable.AlignLeft, Text: res.Code.String()},
		{Align: simpletable.AlignLeft, Text: res.Name},
		{Align: simpletable.AlignLeft, Text: res.Description},
	}
}
