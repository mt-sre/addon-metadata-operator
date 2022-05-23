package list

import (
	"github.com/mt-sre/addon-metadata-operator/cmd/mtcli/list/bundles"
	"github.com/mt-sre/addon-metadata-operator/cmd/mtcli/list/validators"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [command]",
		Short: "Run a list subcommand.",
	}

	cmd.AddCommand(bundles.Cmd())
	cmd.AddCommand(validators.Cmd())

	return cmd
}
