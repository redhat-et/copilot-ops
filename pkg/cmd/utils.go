package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	OpenAI      *openai.Client
	OutputType  string
}

func BuildOpenAIClient(conf Config, nTokens int32, nCompletions int32, engine string) *openai.Client {
	// create OpenAI client
	openAIClient := openai.CreateOpenAIClient(conf.OpenAI.APIKey, conf.OpenAI.OrgID, engine)
	openAIClient.NTokens = nTokens
	openAIClient.NCompletions = nCompletions
	openAIClient.Engine = engine
	return openAIClient
}

func PrepareRequest(cmd *cobra.Command, engine string) (*Request, error) {
	request, _ := cmd.Flags().GetString(FlagRequest)
	write, _ := cmd.Flags().GetBool(FlagWrite)
	path, _ := cmd.Flags().GetString(FlagPath)
	files, _ := cmd.Flags().GetStringArray(FlagFiles)
	if cmd.Name() == CommandEdit {
		file, _ := cmd.Flags().GetString(FlagFiles)
		files = append(files, file)
	}
	filesets, _ := cmd.Flags().GetStringArray(FlagFilesets)
	nTokens, _ := cmd.Flags().GetInt32(FlagNTokens)
	nCompletions, _ := cmd.Flags().GetInt32(FlagNCompletions)
	outputType, _ := cmd.Flags().GetString(FlagOutputType)

	log.Printf("flags:\n")
	log.Printf(" - %-8s: %v\n", FlagRequest, request)
	log.Printf(" - %-8s: %v\n", FlagWrite, write)
	log.Printf(" - %-8s: %v\n", FlagPath, path)
	log.Printf(" - %-8s: %v\n", FlagFiles, files)
	log.Printf(" - %-8s: %v\n", FlagFilesets, filesets)
	log.Printf(" - %-8s: %v\n", FlagNTokens, nTokens)
	log.Printf(" - %-8s: %v\n", FlagNCompletions, nCompletions)
	log.Printf(" - %-8s: %v\n", FlagOutputType, outputType)

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
			// load or ignore failure
			// FIXME: warn when a file has failed to load
			_ = fm.LoadFilesFromGlob(glob)
		}
	}

	if len(filesets) > 0 {
		log.Printf("detected filesets: %v\n", filesets)
		for _, name := range filesets {
			fileset := conf.FindFileset(name)
			if fileset == nil {
				return nil, fmt.Errorf("fileset %s not found in %s", name, ConfigFile)
			}
			for _, glob := range fileset.Files {
				// FIXME: check error here
				_ = fm.LoadFilesFromGlob(glob)
			}
		}
	}

	fm.LogDump()
	filemapText := fm.EncodeToInputText()

	// create OpenAI client
	openAIClient := BuildOpenAIClient(conf, nTokens, nCompletions, engine)
	log.Printf("Model in use: " + openAIClient.Engine)

	r := Request{
		Config:      conf,
		Filemap:     fm,
		FilemapText: filemapText,
		UserRequest: request,
		IsWrite:     write,
		OpenAI:      openAIClient,
		OutputType:  outputType,
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
		fmOutput, err := r.Filemap.EncodeToInputTextFullPaths(r.OutputType)
		if err != nil {
			return err
		}

		stringOut := strings.ReplaceAll(fmOutput, "\\n", "\n")

		log.Printf("\n%s\n", stringOut)
		log.Printf("use --write to actually update files\n")
	}

	return nil
}

// AddRequestFlags Appends flags to the given command which are then used at the command-line.
func AddRequestFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(
		FlagRequest, "r", "",
		"Requested changes in natural language (empty request will surprise you!)",
	)

	cmd.Flags().BoolP(
		FlagWrite, "w", false,
		"Write changes to the repo files (if not set the patch is printed to stdout)",
	)

	cmd.Flags().StringP(
		FlagPath, "p", ".",
		"Path to the root of the repo",
	)

	cmd.Flags().StringP(
		FlagOutputType, "o", "json",
		"How to format output",
	)
}
