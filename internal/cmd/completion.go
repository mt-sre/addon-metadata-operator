package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	mtcli.AddCommand(completionCmd)
}

var (
	completionExample = []string{
		"  # Installing bash completion on Linux",
		"  ## If bash-completion is not installed on Linux, install the 'bash-completion' package",
		"  ## via your distribution's package manager.",
		"  ## Load the mtcli completion code for bash into the current shell",
		"  source <(mtcli completion bash)",
		"  ## Write bash completion code to a file and source it from .bash_profile",
		"  mtcli completion bash > ~/.mtcli/completion.bash.inc",
		`  printf`,
		"  # Mtcli shell completion",
		"  source '$HOME/.mtcli/completion.bash.inc'",
		`  " >> $HOME/.bash_profile"`,
		"  source $HOME/.bash_profile",
		"",
		"  # Load the mtcli completion code for zsh into the current shell",
		"  source <(mtcli completion zsh)",
		"  # Set the mtcli completion code for zsh to autoload on startup",
		"  mtcli completion zsh > ${fpath[1]}/_archsugar",
	}
	completionCmd = &cobra.Command{
		Use:       "completion SHELL",
		Short:     "Output shell completion code for the specified shell (bash or zsh)",
		Example:   strings.Join(completionExample, "\n"),
		Run:       completionMain,
		ValidArgs: []string{"zsh", "bash"},
	}
)

func completionMain(cmd *cobra.Command, args []string) {
	err := cobra.OnlyValidArgs(cmd, args)
	if err != nil {
		log.Fatal(err)
	}

	shell := args[0]
	switch shell {
	case "bash":
		err = mtcli.GenBashCompletion(os.Stdout)
	case "zsh":
		err = mtcli.GenZshCompletion(os.Stdout)
	}

	if err != nil {
		log.Fatalf("Can't generate completion code for shell %v, got %v.\n", shell, err)
	}

}
