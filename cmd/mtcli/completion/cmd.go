package completion

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

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
		"  mtcli completion zsh > ${fpath[1]}/_mtcli",
	}
)

func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:       "completion SHELL",
		Short:     "Output shell completion code for the specified shell (bash or zsh)",
		Example:   strings.Join(completionExample, "\n"),
		RunE:      run,
		ValidArgs: []string{"zsh", "bash"},
	}
}

func run(cmd *cobra.Command, args []string) error {
	if err := cobra.OnlyValidArgs(cmd, args); err != nil {
		return fmt.Errorf("parsing arguments: %w", err)
	}

	shell := args[0]
	switch shell {
	case "bash":
		if err := cmd.Root().GenBashCompletion(os.Stdout); err != nil {
			return fmt.Errorf("generating bash completions: %w", err)
		}
	case "zsh":
		if err := cmd.Root().GenZshCompletion(os.Stdout); err != nil {
			return fmt.Errorf("generating zsh completions: %w", err)
		}
	}

	return nil
}
