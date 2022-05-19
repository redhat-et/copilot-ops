package cmd

import (
	"log"

	"github.com/redhat-et/copilot-ops/pkg/openai"
	"github.com/spf13/cobra"
)

// NewPatchCmd creates the `copilot-ops patch` CLI command
func NewPatchCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use: "patch",

		Short: "Proposes a patch to the repo",

		Long: "Patch takes a request in natural language, packs the related files from the repo, calls AI engine to suggest code changes based on the request, and optionally applies the suggested changes to the repo.",

		Example: `  copilot-ops patch --request 'Add a new secret containing a pre-generated self signed SSL certificate, mount that secret from the syncthing deployment and also the volsync operator deployment, set syncthing configuration to serve with HTTPS using the mounted secret, and add a new go file with a code example that trusts a mounted certificate for the volsync operator pod' --fileset syncthing`,

		RunE: RunPatch,
	}

	AddRequestFlags(cmd)

	return cmd
}

// RunPatch is the implementation of the `copilot-ops patch` command
func RunPatch(cmd *cobra.Command, args []string) error {

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
