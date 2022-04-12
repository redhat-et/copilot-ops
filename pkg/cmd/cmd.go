// cmd contains the commands implementation of the CLI tool
// The main entry point function is cmd.Execute()
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Execute the CLI and exit
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

// RootCmd is the bare `copilot-ops` CLI command which shows the available subcommands
var RootCmd = &cobra.Command{
	Use: "copilot-ops",

	Long: `copilot-ops is a workflow automation tool that proposes an intelligent patches on a repo,
  using natural language AI engines (openai.com codex bring-your-own-token),
  and can be used to implement github bots, editor extensions, and more.
`,

	Example: `  copilot-ops patch --help`,

	// Usage on every error is too noisy and makes it harder
	// to read the error message, so disabling it
	SilenceUsage: true,
}

func init() {
	// Add subcommands of the root command
	RootCmd.AddCommand(PatchCmd)
}
