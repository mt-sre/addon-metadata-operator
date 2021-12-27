package validate

import (
	"github.com/alexeyco/simpletable"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

var (
	statusSuccess = utils.Green("Success")
	statusFailed  = utils.Red("Failed")
	statusError   = utils.IntenselyBoldRed("Error")
)

func getTableHeaders() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "STATUS"},
			{Align: simpletable.AlignCenter, Text: "CODE"},
			{Align: simpletable.AlignCenter, Text: "NAME"},
			{Align: simpletable.AlignCenter, Text: "DESCRIPTION"},
			{Align: simpletable.AlignCenter, Text: "FAILURE MESSAGE"},
		},
	}
}

func newSuccessTableRow(v utils.Validator) []*simpletable.Cell {
	return newTableRow(v, statusSuccess, "")
}

func newFailedTableRow(v utils.Validator, failureMsg string) []*simpletable.Cell {
	return newTableRow(v, statusFailed, failureMsg)
}

func newErrorTableRow(v utils.Validator, err error) []*simpletable.Cell {
	return newTableRow(v, statusError, err.Error())
}

func newTableRow(v utils.Validator, status, failureMsg string) []*simpletable.Cell {
	if failureMsg == "" {
		failureMsg = "None"
	}
	return []*simpletable.Cell{
		{Align: simpletable.AlignLeft, Text: status},
		{Align: simpletable.AlignLeft, Text: v.Code},
		{Align: simpletable.AlignLeft, Text: v.Name},
		{Align: simpletable.AlignLeft, Text: v.Description},
		{Align: simpletable.AlignLeft, Text: failureMsg},
	}
}
