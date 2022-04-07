package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const FLAG_REQUEST = "request"
const FLAG_WRITE = "write"
const FLAG_PATH = "path"
const FLAG_FILES = "file"
const FLAG_FILESETS = "fileset"

// PatchCmd is a command
var PatchCmd = &cobra.Command{
	Use: "patch",

	Short: "Proposes a patch to the repo",

	Long: "Patch takes a request in natural language, packs the related files from the repo, calls AI engine to suggest code changes based on the request, and finally applies the suggested changes to the repo.",

	Example: `  copilot-ops patch --request 'Add a new secret containing a pre-generated self signed SSL certificate, mount that secret from the syncthing deployment and also the volsync operator deployment, set syncthing configuration to serve with HTTPS using the mounted secret, and add a new go file with a code example that trusts a mounted certificate for the volsync operator pod' --fileset syncthing`,

	Run: func(cmd *cobra.Command, args []string) {
		request, _ := cmd.Flags().GetString(FLAG_REQUEST)
		write, _ := cmd.Flags().GetBool(FLAG_WRITE)
		path, _ := cmd.Flags().GetString(FLAG_PATH)
		files, _ := cmd.Flags().GetStringArray(FLAG_FILES)
		filesets, _ := cmd.Flags().GetStringArray(FLAG_FILESETS)

		fmt.Printf("patch request  = %v\n", request)
		fmt.Printf("patch write    = %v\n", write)
		fmt.Printf("patch path     = %v\n", path)
		fmt.Printf("patch files    = %v\n", files)
		fmt.Printf("patch filesets = %v\n", filesets)
	},
}

func init() {
	RootCmd.AddCommand(PatchCmd)

	PatchCmd.Flags().StringP(
		FLAG_REQUEST, "r", "",
		"Requested changes in natural language (empty request will surprise you!)",
	)

	PatchCmd.Flags().BoolP(
		FLAG_WRITE, "w", false,
		"Write changes to the repo files (if not set the patch is printed to stdout)",
	)

	PatchCmd.Flags().StringP(
		FLAG_PATH, "p", ".",
		"Path to the root of the repo",
	)

	PatchCmd.Flags().StringArrayP(
		FLAG_FILES, "f", []string{},
		"File paths (glob) to be considered for the patch (can be specified multiple times)",
	)

	PatchCmd.Flags().StringArrayP(
		FLAG_FILESETS, "s", []string{},
		"Fileset names (defined in .copilot-ops.yaml) to be considered for the patch (can be specified multiple times)",
	)
}
