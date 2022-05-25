// cmd contains the commands implementation of the CLI tool
// The main entry point function is cmd.Execute()
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Execute the CLI and exit
func Execute() {
	err := NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func NewRootCmd() *cobra.Command {
	// the root command shows the available subcommands
	cmd := &cobra.Command{
		Use: "copilot-ops",

		Long: `copilot-ops is a workflow automation tool that proposes an intelligent patches on a repo,
	  using natural language AI engines (openai.com codex bring-your-own-token),
	  and can be used to implement github bots, editor extensions, and more.
	`,

		Example: `  copilot-ops generate --help`,

		// Usage on every error is too noisy and makes it harder
		// to read the error message, so disabling it
		SilenceUsage: true,
	}

	// Add subcommands of the root command
	cmd.AddCommand(NewGenerateCmd())
	cmd.AddCommand(NewEditCmd())

	return cmd
}
