package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(run completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ run completion bash > /etc/bash_completion.d/run
  # macOS:
  $ run completion bash > /usr/local/etc/bash_completion.d/run

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ run completion zsh > "${fpath[1]}/_run"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ run completion fish | source

  # To load completions for each session, execute once:
  $ run completion fish > ~/.config/fish/completions/run.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var target io.Writer
		var targetPath, _ = cmd.Flags().GetString("target")

		if targetPath == "" {
			target = os.Stdout
		} else {
			target, _ = os.OpenFile(targetPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		}

		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(target)
		case "zsh":
			cmd.Root().GenZshCompletion(target)
		case "fish":
			cmd.Root().GenFishCompletion(target, true)
		}
	},
}

func init() {
	completionCmd.Flags().String("target", "", "target file")
	rootCmd.AddCommand(completionCmd)
}
