package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/mt-sre/addon-metadata-operator/cmd/mtcli/bundle"
	"github.com/mt-sre/addon-metadata-operator/cmd/mtcli/completion"
	"github.com/mt-sre/addon-metadata-operator/cmd/mtcli/list"
	"github.com/mt-sre/addon-metadata-operator/cmd/mtcli/validate"
	"github.com/mt-sre/addon-metadata-operator/cmd/mtcli/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var verbose bool

func main() {
	code := 0

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer func() {
		stop()

		os.Exit(code)
	}()

	rootCmd := generateRootCmd()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stdout, err)

		code = 1
	}
}

func generateRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "mtcli",
		Short: "Managed Tenants CLI swiss army knife.",
	}

	rootCmd.AddCommand(bundle.Cmd())
	rootCmd.AddCommand(completion.Cmd())
	rootCmd.AddCommand(list.Cmd())
	rootCmd.AddCommand(validate.Cmd())
	rootCmd.AddCommand(version.Cmd())

	flags := rootCmd.PersistentFlags()
	flags.BoolVarP(
		&verbose,
		"verbose",
		"v",
		verbose,
		"verbose output",
	)

	cobra.OnInitialize(setLogFormatter, setLogLevel)

	return rootCmd
}

func setLogFormatter() {
	formatter := &log.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	}
	log.SetFormatter(formatter)
}

func setLogLevel() {
	log.SetLevel(log.InfoLevel)
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}
