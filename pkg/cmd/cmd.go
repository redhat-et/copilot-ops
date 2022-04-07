package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Execute the CLI
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// RootCmd is the root command of the CLI
var RootCmd = &cobra.Command{
	Use: "copilot-ops",

	Long: `copilot-ops is a workflow automation tool that proposes an intelligent patches on a repo,
  using natural language AI engines (openai.com codex bring-your-own-token),
  and can be used to implement github bots, editor extensions, and more.
`,

	Example: `  copilot-ops patch --help`,
}
