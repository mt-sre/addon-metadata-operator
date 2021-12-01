package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	mtcli.AddCommand(versionCmd)
}

// Injected by goreleaser through ldflags (see .goreleaser.yml)
var (
	version = "0.0.0"
	commit  = "local"
	builtBy = "local"
	date    = "local"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show mtcli version information.",
		Args:  cobra.NoArgs,
		Run:   versionMain,
	}
)

func versionMain(cmd *cobra.Command, args []string) {
	fmt.Printf("mtcli version: %v, commit: %v, builtBy: %v (%v)\n", version, commit, builtBy, date)
}
