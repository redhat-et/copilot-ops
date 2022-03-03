/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// suggestCmd represents the suggest command
var suggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Generates YAMLs for the provided issue and given files, if any",
	Long: `Given an input prompt and some optional files, this command
suggests new YAML files to address the issue.

If --dry-run is provided, the suggested YAMLs are printed to stdout.
`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(suggestCmd)

	// Here you will define your flags and configuration settings.
	suggestCmd.PersistentFlags().StringP("path", "p", ".", "path to the root of the git repo")
	suggestCmd.PersistentFlags().StringP("issue", "i", "", "description of required changes. Files must be referenced by their '@filetag'")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// suggestCmd.PersistentFlags().String("foo", "", "A help for foo")

	suggestCmd.PersistentFlags().StringArrayP("files", "f", []string{}, "files to be modified")
	suggestCmd.Flags().Bool("dry-run", false, "suggested files are printed to stdout")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// suggestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// mark commands as required
	suggestCmd.MarkFlagRequired("issue")
}
