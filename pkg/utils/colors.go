package utils

import "github.com/fatih/color"

var (
	// colors for stdout
	Green  = color.New(color.FgGreen).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Bold   = color.New(color.Bold).SprintFunc()
)
