package bundle

import (
	"github.com/mt-sre/addon-metadata-operator/cmd/mtcli/bundle/validate"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bundle [command]",
		Short: "Run a bundle subcommand.",
	}

	cmd.AddCommand(validate.Cmd())

	return cmd
}
