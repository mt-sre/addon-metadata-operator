package version

import (
	"fmt"
	"os"

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
	fmt.Fprintln(os.Stdout, cli.Version())

	return nil
}
