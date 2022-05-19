package cmd

import (
	"fmt"
	"log"

	"github.com/redhat-et/copilot-ops/pkg/filemap"
	"github.com/redhat-et/copilot-ops/pkg/openai"
	"github.com/spf13/cobra"
)

// NewGenerateCmd creates the `copilot-ops patch` CLI command
func NewGenerateCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use: "generate",

		Short: "Proposes a new files to the repo",

		Long: "Generate takes a request in natural language, packs the related files from the repo, calls AI engine to suggest generating new code based on the request, and optionally applies the suggested changes to the repo.",

		Example: `  copilot-ops patch --request 'Add a new secret containing a pre-generated self signed SSL certificate, mount that secret from the syncthing deployment and also the volsync operator deployment, set syncthing configuration to serve with HTTPS using the mounted secret, and add a new go file with a code example that trusts a mounted certificate for the volsync operator pod' --fileset syncthing`,

		RunE: RunGenerate,
	}

	AddRequestFlags(cmd)

	// cmd.Flags().Int16P()

	return cmd
}

// RunGenerate is the implementation of the `copilot-ops patch` command
func RunGenerate(cmd *cobra.Command, args []string) error {

	r, err := PrepareRequest(cmd, openai.OpenAICodeDavinciV2)
	if err != nil {
		return err
	}

	input := PrepareGenerateInput(r.UserRequest, r.FilemapText)
	log.Printf("requesting generate from OpenAI: %s", input)

	output, err := r.OpenAI.GenerateCode(input)
	if err != nil {
		return err
	}
	log.Printf("received output from OpenAI: %q\n", output)

	// err = r.Filemap.DecodeFromOutput(output)
	r.Filemap = filemap.NewFilemap()
	r.Filemap.DecodeFromOutput(output)
	if err != nil {
		// TODO try other way to decode the output to a fileset

		// fallback - generate a new filename and put the content inside
		newFileName := "generated-by-copilot-ops"
		r.Filemap.Files = map[string]filemap.File{
			newFileName: {Path: newFileName, Content: output, Tag: newFileName},
		}
	}

	return PrintOrWriteOut(r)
}

func PrepareGenerateInput(userInput string, encodedFiles string) string {

	return fmt.Sprintf(`
## This document contains:
## 1. Instructions describing new files that need to be created
## 2. The existing files, each separated by a '%s'
## 3. The newly created files, which are separated by '%s', and terminated by a '%s' sequence 

## 1. Instructions for the new files:
%s

## 2. The existing files:
%s

## 3. The newly-created files:
`, filemap.FileDelimeter, filemap.FileDelimeter, openai.CompletionEndOfSequence, userInput, encodedFiles)
}
