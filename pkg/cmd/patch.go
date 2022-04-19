package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/redhat-et/copilot-ops/pkg/filemap"
	"github.com/redhat-et/copilot-ops/pkg/openai"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// PatchCmd is the `copilot-ops patch` CLI command
var PatchCmd = &cobra.Command{
	Use: "patch",

	Short: "Proposes a patch to the repo",

	Long: "Patch takes a request in natural language, packs the related files from the repo, calls AI engine to suggest code changes based on the request, and finally applies the suggested changes to the repo.",

	Example: `  copilot-ops patch --request 'Add a new secret containing a pre-generated self signed SSL certificate, mount that secret from the syncthing deployment and also the volsync operator deployment, set syncthing configuration to serve with HTTPS using the mounted secret, and add a new go file with a code example that trusts a mounted certificate for the volsync operator pod' --fileset syncthing`,

	RunE: RunPatch,
}

func init() {
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
		"Fileset names (defined in "+CONFIG_FILE+") to be considered for the patch (can be specified multiple times)",
	)
}

const FLAG_REQUEST = "request"
const FLAG_WRITE = "write"
const FLAG_PATH = "path"
const FLAG_FILES = "file"
const FLAG_FILESETS = "fileset"

// RunPatch is the implementation of the `copilot-ops patch` command
func RunPatch(cmd *cobra.Command, args []string) error {

	request, _ := cmd.Flags().GetString(FLAG_REQUEST)
	write, _ := cmd.Flags().GetBool(FLAG_WRITE)
	path, _ := cmd.Flags().GetString(FLAG_PATH)
	files, _ := cmd.Flags().GetStringArray(FLAG_FILES)
	filesets, _ := cmd.Flags().GetStringArray(FLAG_FILESETS)

	log.Printf("flags:\n")
	log.Printf(" - %-8s: %v\n", FLAG_REQUEST, request)
	log.Printf(" - %-8s: %v\n", FLAG_WRITE, write)
	log.Printf(" - %-8s: %v\n", FLAG_PATH, path)
	log.Printf(" - %-8s: %v\n", FLAG_FILES, files)
	log.Printf(" - %-8s: %v\n", FLAG_FILESETS, filesets)

	// Handle --path by changing the working directory
	// so that every file name we refer to is relative to path
	if path != "" {
		if err := os.Chdir(path); err != nil {
			return err
		}
	}

	// Load the config from file if it exists, but if it doesn't exist
	// we'll just use the defaults and continue without error.
	// Errors here might return if the file exists but is invalid.
	config := Config{}
	err := config.LoadFile(CONFIG_FILE)
	if err != nil {
		return err
	}
	log.Printf("config: %+v\n", config)

	fm := filemap.NewFilemap()

	if len(files) > 0 {
		for _, glob := range files {
			fm.LoadFilesFromGlob(glob)
		}
	}

	if len(filesets) > 0 {
		for _, name := range filesets {
			fileset := config.GetFileset(name)
			if fileset == nil {
				return fmt.Errorf("fileset %s not found in %s", name, CONFIG_FILE)
			}
			for _, glob := range fileset.Files {
				fm.LoadFilesFromGlob(glob)
			}
		}
	}

	fm.LogDump()

	input, err := fm.EncodeToInputText()
	if err != nil {
		return err
	}

	ai, err := openai.CreateOpenAIClient()
	if err != nil {
		return err
	}

	output, err := ai.EditCode(input, request)
	if err != nil {
		return err
	}

	fmOut := filemap.NewFilemap()
	err = fmOut.DecodeFromOutput(output)
	if err != nil {
		return err
	}

	fmOut.LogDump()

	if write {
		log.Printf("updating files ...\n")
		err = fmOut.WriteUpdatesToFiles()
		if err != nil {
			return err
		}
	} else {
		log.Printf("use --write to actually update files\n")
	}

	return nil
}

const CONFIG_FILE = ".copilot-ops.yaml"

type Config struct {
	Filesets []ConfigFilesets
}

type ConfigFilesets struct {
	Name  string
	Files []string
}

func (c *Config) LoadFile(path string) error {
	text, _ := os.ReadFile(path)
	return yaml.Unmarshal(text, c)
}

func (c *Config) GetFileset(name string) *ConfigFilesets {
	for _, fileset := range c.Filesets {
		if fileset.Name == name {
			return &fileset
		}
	}
	return nil
}
