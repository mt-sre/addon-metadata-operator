package version

import (
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/internal/cli"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show mtcli version information.",
		Args:  cobra.NoArgs,
		RunE:  run,
	}
}

func run(cmd *cobra.Command, args []string) error {
	fmt.Fprintln(cmd.OutOrStdout(), cli.Version())

	return nil
}
