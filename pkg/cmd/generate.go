package cmd

import (
	"context"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/redhat-et/copilot-ops/pkg/cmd/config"
	"github.com/redhat-et/copilot-ops/pkg/filemap"
	"github.com/redhat-et/copilot-ops/pkg/openai"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/spf13/cobra"
)

// NewGenerateCmd creates the `copilot-ops patch` CLI command.
func NewGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: CommandGenerate,

		Short: "Proposes a new files to the repo",

		Long: "Generate takes a request in natural language, packs the related " +
			"files from the repo, calls AI engine to suggest generating new code " +
			"based on the request, and optionally applies the suggested changes to the repo.",

		Example: `  copilot-ops generate --file examples/app1/mysql-pvc.yaml --request` +
			`'Generate a pod that mounts the PVC. Set the pod resources requests and ` +
			`limits to 4 cpus and 5 Gig of memory.'`,

		RunE: RunGenerate,
	}

	AddRequestFlags(cmd)

	// generate-specific flags
	cmd.Flags().StringArrayP(
		FlagFilesFull, FlagFilesShort, []string{},
		"File paths (glob) to be considered for the patch (can be specified multiple times)",
	)

	cmd.Flags().StringArrayP(
		FlagFilesetsFull, FlagFilesetsShort, []string{},
		"Fileset names (defined in "+config.ConfigFile+") to be considered for the patch (can be specified multiple times)",
	)

	cmd.Flags().Int32P(
		FlagNTokensFull, FlagNTokensShort, DefaultTokens,
		"Max number of tokens to generate",
	)

	cmd.Flags().Int32P(
		FlagNCompletionsFull, FlagNCompletionsShort, DefaultCompletions,
		"Number of completions to generate",
	)

	return cmd
}

// RunGenerate is the implementation of the `copilot-ops generate` command.
func RunGenerate(cmd *cobra.Command, args []string) error {
	r, err := PrepareRequest(cmd, openai.OpenAICodeDavinciV2)
	if err != nil {
		return err
	}
	input := PrepareGenerateInput(r.UserRequest, r.FilemapText)
	log.Printf("requesting generate from OpenAI: %s", input)

	// generate a response from OpenAI
	model := openai.OpenAICodeDavinciV2
	resp, err := r.OpenAI.CreateCompletion(
		context.TODO(),
		openai.OpenAICodeDavinciV2,
		gogpt.CompletionRequest{
			Model:       &model,
			Prompt:      input,
			MaxTokens:   int(r.NTokens),
			N:           int(r.NCompletions),
			Temperature: 0.0,
			Stop:        []string{openai.CompletionEndOfSequence},
		},
	)
	if err != nil {
		return fmt.Errorf("got error from OpenAI: %w", err)
	}

	// decode the response
	r.Filemap = filemap.NewFilemap()
	log.Printf("decoding output")
	for _, s := range resp.Choices {
		if err == nil {
			err = r.Filemap.DecodeFromOutput(s.Text)
		}
	}

	if err == nil {
		return PrintOrWriteOut(r)
	}

	// HACK: try other way to decode the output to a fileset
	log.Printf("decoding failed, got error: %s", err)
	// fallback - generate new files and put the content inside
	generateNewFiles(r, resp.Choices)

	return PrintOrWriteOut(r)
}

// PrepareGenerateInput Accepts the userInput and all of the files encoded as a string,
// and formats them as a prompt to be sent off to OpenAI.
func PrepareGenerateInput(userInput string, encodedFiles string) string {
	// HACK: prompt wording needs to be adjusted to improve accuracy
	var prompt = ""
	var withFiles = len(encodedFiles) > 0

	// preamble
	prompt += preamble(withFiles)

	// instructions
	prompt += instructions(withFiles)

	// prompt the AI for a response
	prompt += callToActionSequence(userInput, encodedFiles)
	return prompt
}

// preamble Returns the preamble for the generation prompt, with varied text
// depending on whether or not the prompt will be including other relevant YAML
// files.
func preamble(withFiles bool) string {
	if withFiles {
		return `## This document contains instructions for a new Kubernetes YAML that needs to be created,
## along with the relevant YAMLs for context, and the resultant YAML.`
	}
	return `## This document contains instructions for a new Kubernetes YAML that needs to be created,
## and the resultant YAML.`
}

// instructions Returns the sequence in the prompt which details the ordering of the
// document for the AI, and what it should expect when parsing the tokens.
func instructions(withFiles bool) string {
	var numInstructions int8 = 1

	// instructions
	prompt := fmt.Sprintf(`
##
## The structure of the document is as follows:
## %d. Description of the desired YAML`, numInstructions)
	numInstructions++

	// mention that extra YAMLs will be provided for context
	if withFiles {
		prompt += fmt.Sprintf(`
## %d. The existing YAMLs, each separated by a '%s'`, numInstructions, filemap.FileDelimeter)
		numInstructions++
	}

	// instruction for the generated code
	prompt += fmt.Sprintf(`
## %d. The new YAML, terminated by an '%s'`, numInstructions, openai.CompletionEndOfSequence)
	prompt += "\n"

	return prompt
}

// callToActionSequence Creates the section which includes the actual request
// for the generated YAML, along with the encodedFiles for context if those are also needed.
func callToActionSequence(request string, encodedFiles string) string {
	// reset counter
	numInstructions := 1

	// add the user input
	prompt := fmt.Sprintf(`
## %d. Instructions for the new Kubernetes YAML:
%s
`, numInstructions, request)
	numInstructions++

	// add the encoded files if they exist
	if strings.TrimSpace(encodedFiles) != "" {
		prompt += fmt.Sprintf(`
## %d. Existing YAMLs:
%s
`, numInstructions, encodedFiles)
		numInstructions++
	}

	// add the completion sequence
	prompt += fmt.Sprintf(`
## %d. The new YAML:
`, numInstructions)
	return prompt
}

// generateNewFiles Creates a new file for every requested completion,
// and stores them in the "generated-by-copilot-ops" directory.
func generateNewFiles(req *Request, sepOutput []gogpt.CompletionChoice) {
	newMap := make(map[string]filemap.File)
	for i := 0; i < len(sepOutput); i++ {
		// set file name + path here
		newFileName := "generated-by-copilot-ops" + fmt.Sprint(i+1) + ".yaml"
		newFilePath := path.Join("generated-by-copilot-ops", newFileName)

		var newFile filemap.File
		newFile.Content = sepOutput[i].Text
		newFile.Path = newFilePath
		newFile.Name = newFileName

		newMap[newFilePath] = newFile
	}

	req.Filemap.Files = newMap
}
