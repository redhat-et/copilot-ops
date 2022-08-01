package cmd

import (
	"context"
	"fmt"

	"github.com/redhat-et/copilot-ops/pkg/filemap"
	"github.com/redhat-et/copilot-ops/pkg/openai"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/spf13/cobra"
)

// NewEditCmd Creates the `copilot-ops edit` CLI command.
func NewEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: CommandEdit,

		Short: "Edits a single file provided to the CLI",

		Long: "Given a file and a request, edit will pinpoint what changes are being requested and attempt to make them.",

		Example: `  copilot-ops edit --file examples/app1/mysql-pvc.yaml --request 'Increase the size of the PVC to 100Gi'`,

		RunE: RunEdit,
	}

	AddRequestFlags(cmd)

	// flag to add a file
	cmd.Flags().StringP(
		FlagFilesFull, FlagFilesShort, "",
		"File path to the document which should be edited.",
	)

	return cmd
}

// RunEdit Runs when the `edit` command is invoked.
func RunEdit(cmd *cobra.Command, args []string) error {
	r, err := PrepareRequest(cmd, openai.OpenAICodeDavinciEditV1)
	if err != nil {
		return err
	}

	// trigger GPT-3 to preserve the @tagname format in the file
	editSuffix := fmt.Sprintf("The resulting file should preserve the '# %stagname'"+
		" format used to identify the YAML(s).", filemap.FileTagPrefix)
	editInstruction := fmt.Sprintf("%s\n\n%s", r.UserRequest, editSuffix)

	// create a client for editing
	model := openai.OpenAICodeDavinciEditV1
	response, err := r.OpenAI.Edits(
		context.TODO(),
		gogpt.EditsRequest{
			Model:       &model,
			Input:       r.FilemapText,
			Instruction: editInstruction,
			// FIXME: edit more than one file
			N:           1,
			Temperature: 0.0,
		},
	)
	if err != nil {
		return err
	}
	output := response.Choices[0].Text
	err = r.Filemap.DecodeFromOutput(output)
	if err != nil {
		return err
	}

	return PrintOrWriteOut(r)
}
