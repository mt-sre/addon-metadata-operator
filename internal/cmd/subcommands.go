package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	// Register all subcommands
	mtcli.AddCommand(bundleCmd)
	mtcli.AddCommand(listCmd)
}

var (
	bundleCmd = &cobra.Command{
		Use:   "bundle [command]",
		Short: "Run a bundle subcommand.",
		Run:   func(cmd *cobra.Command, args []string) { _ = cmd.Help() },
	}
	listCmd = &cobra.Command{
		Use:   "list [command]",
		Short: "Run a list subcommand.",
		Run:   func(cmd *cobra.Command, args []string) { _ = cmd.Help() },
	}
)
