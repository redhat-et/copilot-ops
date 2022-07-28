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
	OpenAIURL   string
}

// BuildOpenAIClient Creates and configures an OpenAI client based on the given parameters.
func BuildOpenAIClient(conf Config, nTokens int32, nCompletions int32, engine string, openAIURL string) *openai.Client {
	// create OpenAI client
	openAIClient := openai.CreateOpenAIClient()
	openAIClient.NTokens = nTokens
	openAIClient.NCompletions = nCompletions
	openAIClient.Engine = engine
	// allow override of OpenAI URL for testing
	openAIClient.APIUrl = openAIURL + openai.OpenAIEndpointV1
	return openAIClient
}

// PrepareRequest Processes the user input along with provided environment variables,
// creating a Request object which is used for context in further requests.
func PrepareRequest(cmd *cobra.Command, engine string) (*Request, error) {
	request, _ := cmd.Flags().GetString(FlagRequestFull)
	write, _ := cmd.Flags().GetBool(FlagWriteFull)
	path, _ := cmd.Flags().GetString(FlagPathFull)
	files, _ := cmd.Flags().GetStringArray(FlagFilesFull)
	if cmd.Name() == CommandEdit {
		file, _ := cmd.Flags().GetString(FlagFilesFull)
		files = append(files, file)
	}
	filesets, _ := cmd.Flags().GetStringArray(FlagFilesetsFull)
	nTokens, _ := cmd.Flags().GetInt32(FlagNTokensFull)
	nCompletions, _ := cmd.Flags().GetInt32(FlagNCompletionsFull)
	outputType, _ := cmd.Flags().GetString(FlagOutputTypeFull)
	openAIURL, _ := cmd.Flags().GetString(FlagOpenAIURLFull)

	log.Printf("flags:\n")
	log.Printf(" - %-8s: %v\n", FlagRequestFull, request)
	log.Printf(" - %-8s: %v\n", FlagWriteFull, write)
	log.Printf(" - %-8s: %v\n", FlagPathFull, path)
	log.Printf(" - %-8s: %v\n", FlagFilesFull, files)
	log.Printf(" - %-8s: %v\n", FlagFilesetsFull, filesets)
	log.Printf(" - %-8s: %v\n", FlagNTokensFull, nTokens)
	log.Printf(" - %-8s: %v\n", FlagNCompletionsFull, nCompletions)
	log.Printf(" - %-8s: %v\n", FlagOutputTypeFull, outputType)

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
	_ = conf.Load()

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
	openAIClient := BuildOpenAIClient(conf, nTokens, nCompletions, engine, openAIURL)
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
	if r.IsWrite {
		err := r.Filemap.WriteUpdatesToFiles()
		if err != nil {
			return err
		}
		return nil
	}

	// TODO: print as redirectable / pipeable write stream
	fmOutput, err := r.Filemap.EncodeToInputTextFullPaths(r.OutputType)
	if err != nil {
		return err
	}
	stringOut := strings.ReplaceAll(fmOutput, "\\n", "\n")
	log.Printf("\n%s\n", stringOut)

	return nil
}

// AddRequestFlags Appends flags to the given command which are then used at the command-line.
func AddRequestFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(
		FlagRequestFull, FlagRequestShort, "",
		"Requested changes in natural language (empty request will surprise you!)",
	)

	cmd.Flags().BoolP(
		FlagWriteFull, FlagWriteShort, false,
		"Write changes to the repo files (if not set the patch is printed to stdout)",
	)

	cmd.Flags().StringP(
		FlagPathFull, FlagPathShort, ".",
		"Path to the root of the repo",
	)

	cmd.Flags().StringP(
		FlagOutputTypeFull, FlagOutputTypeShort, "json",
		"How to format output",
	)

	_ = cmd.Flags().StringP(
		FlagOpenAIURLFull,
		FlagOpenAIURLShort,
		openai.OpenAIURL,
		"Domain of the OpenAI API",
	)
}
