package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

func init() {
	mtcli.PersistentFlags().BoolVarP(&mtcliVerbose, "verbose", "v", mtcliVerbose, "verbose output")
	cobra.OnInitialize(setLogFormatter, setLogLevel)
}

var (
	mtcliVerbose = false
	mtcli        = &cobra.Command{
		Use:   "mtcli",
		Short: "Managed Tenants CLI swiss army knife.",
		Run:   mtcliMain,
	}
)

func mtcliMain(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}

// Execute - called by main.main()
func Execute() {
	if err := mtcli.Execute(); err != nil {
		log.Fatalln(err)
	}
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
	if mtcliVerbose {
		log.SetLevel(log.DebugLevel)
	}
}
