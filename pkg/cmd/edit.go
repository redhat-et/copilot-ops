package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/redhat-et/copilot-ops/pkg/filemap"
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

	// trigger GPT-3 to preserve the @tagname format in the file
	editSuffix := fmt.Sprintf("The resulting file should preserve the '# %stagname' format used to identify the YAML(s).", filemap.FILE_TAG_PREFIX)
	editInstruction := fmt.Sprintf("%s\n\n%s", r.UserRequest, editSuffix)

	output, err := r.OpenAI.EditCode(r.FilemapText, editInstruction)
	if err != nil {
		return err
	}

	stringOut := strings.ReplaceAll(string(output), "\\n", "\n")

	log.Printf("received patch from OpenAI: \n%s\n", stringOut)

	err = r.Filemap.DecodeFromOutput(output)
	if err != nil {
		return err
	}

	return PrintOrWriteOut(r)
}
