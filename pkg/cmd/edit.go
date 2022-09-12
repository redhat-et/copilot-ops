package cmd

import (
	"fmt"

	"github.com/redhat-et/copilot-ops/pkg/ai"
	"github.com/redhat-et/copilot-ops/pkg/ai/gpt3"
	"github.com/redhat-et/copilot-ops/pkg/filemap"
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
	r, err := PrepareRequest(cmd)
	if err != nil {
		return err
	}

	// trigger GPT-3 to preserve the @tagname format in the file
	editSuffix := fmt.Sprintf("The resulting file should preserve the '# %stagname'"+
		" format used to identify the YAML(s).", filemap.FileTagPrefix)
	editInstruction := fmt.Sprintf("%s\n\n%s", r.UserRequest, editSuffix)

	// create a client for editing
	client, err := PrepareEditClient(r, r.FilemapText, editInstruction)
	if err != nil {
		return fmt.Errorf("could not create client: %w", err)
	}

	responses, err := client.Edit()
	if err != nil {
		return fmt.Errorf("could not edit files: %w", err)
	}
	output := responses[0]
	err = r.Filemap.DecodeFromOutput(output)
	if err != nil {
		return err
	}

	return PrintOrWriteOut(r)
}

// PrepareEditClient Returns an AI Client which implements the EditClient interface.
func PrepareEditClient(r *Request, input, instruction string) (ai.EditClient, error) {
	var client ai.EditClient

	switch r.Backend {
	case ai.GPT3:
		config := r.Config.OpenAI
		if config == nil {
			return nil, fmt.Errorf("no openai config provided")
		}
		client = gpt3.CreateGPT3EditClient(gpt3.OpenAIConfig{
			Token:   config.APIKey,
			OrgID:   &config.OrgID,
			BaseURL: config.URL,
			// FIXME: edit more than one file
		}, input, instruction, 1, nil, nil)
	case ai.GPTJ:
		return nil, fmt.Errorf("editing is not implemented for gpt-j")
	case ai.BLOOM:
		return nil, fmt.Errorf("editing is not implemented for bloom")
	case ai.OPT:
		return nil, fmt.Errorf("editing is not implemented for opt")
	case ai.Unselected:
		return nil, fmt.Errorf("no backend selected")
	default:
		return nil, fmt.Errorf("selected backend does not implement the edit client")
	}
	return client, nil
}
