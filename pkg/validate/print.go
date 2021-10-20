package validate

import (
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

var (
	equalLines = strings.Repeat("=", 6)

	success = utils.Green("Success")
	failed  = utils.Red("Failed")
	err     = utils.IntenselyBoldRed("Error")
)

func printMetaHeading() {
	fmt.Printf("\n%sRUNNING METADATA VALIDATORS%s\n\n", equalLines, equalLines)
}

func printSuccessMessage(msg string) {
	fmt.Printf("\r%s\t\t%s", msg, success)
}

func printFailureMessage(msg string) {
	fmt.Printf("\r%s\t\t%s", msg, failed)
}

func printErrorMessage(msg string) {
	fmt.Printf("\r%s\t\t%s", msg, err)
}