package cmd

import (
	"log"

	"github.com/redhat-et/copilot-ops/pkg/openai"
	"github.com/spf13/cobra"
)

// NewEditCmd Creates the `copilot-ops edit` CLI command.
func NewEditCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "edit",
		Short:   "Edits a single file provided to the CLI",
		Long:    "Given a file and a prompt, edit will pinpoint what changes are being requested and attempt to make them.",
		Example: `  copilot-ops edit --file mysql-pvc.yaml --request 'Increase the size of the PVC to 100Gi'`,
		RunE:    RunEdit,
	}

	AddRequestFlags(cmd)

	// flag to add a file
	cmd.Flags().StringP(
		FLAG_FILES, "f", "",
		"File path to the document which should be edited.",
	)

	return cmd
}

func RunEdit(cmd *cobra.Command, args []string) error {

	r, err := PrepareRequest(cmd, openai.OpenAICodeDavinciEditV1)
	if err != nil {
		return err
	}

	// TODO continue this edit command
	// parseYAML() -> send to openai
	// response() -> determine which lines to edit
	log.Printf("TODO EDIT: %v\n", r)

	return nil
}
