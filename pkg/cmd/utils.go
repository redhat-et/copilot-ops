package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/redhat-et/copilot-ops/pkg/filemap"
	"github.com/redhat-et/copilot-ops/pkg/openai"
	"github.com/spf13/cobra"
)

type Request struct {
	Config      Config
	Fileset     *ConfigFilesets
	Filemap     *filemap.Filemap
	FilemapText string
	UserRequest string
	IsWrite     bool
	OpenAI      *openai.OpenAIClient
}

func PrepareRequest(cmd *cobra.Command, engine string) (*Request, error) {

	request, _ := cmd.Flags().GetString(FLAG_REQUEST)
	write, _ := cmd.Flags().GetBool(FLAG_WRITE)
	path, _ := cmd.Flags().GetString(FLAG_PATH)
	files, _ := cmd.Flags().GetStringArray(FLAG_FILES)
	if cmd.Name() == COMMAND_EDIT {
		file, _ := cmd.Flags().GetString(FLAG_FILES)
		files = append(files, file)
	}
	filesets, _ := cmd.Flags().GetStringArray(FLAG_FILESETS)
	nTokens, _ := cmd.Flags().GetInt32(FLAG_NTOKENS)

	log.Printf("flags:\n")
	log.Printf(" - %-8s: %v\n", FLAG_REQUEST, request)
	log.Printf(" - %-8s: %v\n", FLAG_WRITE, write)
	log.Printf(" - %-8s: %v\n", FLAG_PATH, path)
	log.Printf(" - %-8s: %v\n", FLAG_FILES, files)
	log.Printf(" - %-8s: %v\n", FLAG_FILESETS, filesets)
	log.Printf(" - %-8s: %v\n", FLAG_NTOKENS, nTokens)

	// Handle --path by changing the working directory
	// so that every file name we refer to is relative to path
	if path != "" {
		if err := os.Chdir(path); err != nil {
			return nil, err
		}
	}

	// Load the config from file if it exists, but if it doesn't exist
	// we'll just use the defaults and continue without error.
	// Errors here might return if the file exists but is invalid.
	conf := Config{}
	err := conf.Load()
	if err != nil {
		return nil, err
	}

	// WARNING we should not consider printing conf with its secret keys
	log.Printf("Filesets: %+v\n", conf.Filesets)

	fm := filemap.NewFilemap()

	if len(files) > 0 {
		log.Printf("loading files from command line: %v\n", files)
		for _, glob := range files {
			fm.LoadFilesFromGlob(glob)
		}
	}

	if len(filesets) > 0 {
		log.Printf("detected filesets: %v\n", filesets)
		for _, name := range filesets {
			fileset := conf.FindFileset(name)
			if fileset == nil {
				return nil, fmt.Errorf("fileset %s not found in %s", name, CONFIG_FILE)
			}
			for _, glob := range fileset.Files {
				fm.LoadFilesFromGlob(glob)
			}
		}
	}

	fm.LogDump()

	filemapText, err := fm.EncodeToInputText()
	if err != nil {
		return nil, err
	}
	log.Printf("decoded input: %q\n", filemapText)

	// create OpenAI client
	openAIClient := openai.CreateOpenAIClient(conf.OpenAI.ApiKey, conf.OpenAI.OrgId, engine)
	openAIClient.NTokens = nTokens

	r := Request{
		Config:      conf,
		Filemap:     fm,
		FilemapText: filemapText,
		UserRequest: request,
		IsWrite:     write,
		OpenAI:      openAIClient,
	}

	return &r, nil
}

// PrintOrWriteOut Accepts a request object and writes the contents of the filemap
// to the disk if specified, otherwise it prints to STDOUT.
func PrintOrWriteOut(r *Request) error {

	// dump the state of the FileMap
	r.Filemap.LogDump()

	if r.IsWrite {
		log.Printf("updating files ...\n")
		err := r.Filemap.WriteUpdatesToFiles()
		if err != nil {
			return err
		}
	} else {
		// TODO: Add output formatting control
		// just encode the output and print it to stdout
		// TODO: print as redirectable / pipeable write stream
		fmOutput, err := r.Filemap.EncodeToInputTextFullPaths()
		if err != nil {
			return err
		}
		log.Printf("\n%s\n", fmOutput)
		log.Printf("use --write to actually update files\n")
	}

	return nil
}

// AddRequestFlags Appends flags to the given command which are then used at the command-line.
func AddRequestFlags(cmd *cobra.Command) {

	cmd.Flags().StringP(
		FLAG_REQUEST, "r", "",
		"Requested changes in natural language (empty request will surprise you!)",
	)

	cmd.Flags().BoolP(
		FLAG_WRITE, "w", false,
		"Write changes to the repo files (if not set the patch is printed to stdout)",
	)

	cmd.Flags().StringP(
		FLAG_PATH, "p", ".",
		"Path to the root of the repo",
	)

}
