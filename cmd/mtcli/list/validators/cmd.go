package validators

import (
	"fmt"
	"os"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/internal/cli"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	_ "github.com/mt-sre/addon-metadata-operator/pkg/validator/register"
	"github.com/spf13/cobra"
)

func examples() string {
	return strings.Join([]string{
		"  # List all the registered validators.",
		"  mtcli list validators",
	}, "\n")
}

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:     "validators",
		Short:   "List all the registered validators.",
		Example: examples(),
		RunE:    run,
	}
}

func run(cmd *cobra.Command, args []string) error {
	runner, err := validator.NewRunner()
	if err != nil {
		return fmt.Errorf("listing validators: %s\n", err)
	}

	table, err := cli.NewTable(
		cli.WithHeaders{"CODE", "NAME", "DESCRIPTION"},
	)
	if err != nil {
		return fmt.Errorf("initializing table: %w", err)
	}

	for _, v := range runner.GetValidators() {
		table.WriteRow(cli.TableRow{
			cli.Field{Value: v.Code().String()},
			cli.Field{Value: v.Name()},
			cli.Field{Value: v.Description()},
		})
	}

	fmt.Fprintln(os.Stdout, table.String())
	fmt.Fprintln(os.Stdout)

	return nil
}
