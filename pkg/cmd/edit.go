package cmd

import (
	"log"

	"github.com/redhat-et/copilot-ops/pkg/openai"
	"github.com/spf13/cobra"
)

// NewEditCmd Creates the `copilot-ops edit` CLI command.
func NewEditCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use: COMMAND_EDIT,

		Short: "Edits a single file provided to the CLI",

		Long: "Given a file and a request, edit will pinpoint what changes are being requested and attempt to make them.",

		Example: `  copilot-ops edit --file examples/app1/mysql-pvc.yaml --request 'Increase the size of the PVC to 100Gi'`,

		RunE: RunEdit,
	}

	AddRequestFlags(cmd)

	// flag to add a file
	cmd.Flags().StringP(
		FLAG_FILES, "f", "",
		"File path to the document which should be edited.",
	)

	return cmd
}

// RunEdit Runs when the `edit` command is
func RunEdit(cmd *cobra.Command, args []string) error {

	r, err := PrepareRequest(cmd, openai.OpenAICodeDavinciEditV1)
	if err != nil {
		return err
	}

	output, err := r.OpenAI.EditCode(r.FilemapText, r.UserRequest)
	if err != nil {
		return err
	}

	log.Printf("received patch from OpenAI: %q\n", output)

	err = r.Filemap.DecodeFromOutput(output)
	if err != nil {
		return err
	}

	return PrintOrWriteOut(r)
}
