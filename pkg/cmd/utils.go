package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/redhat-et/copilot-ops/pkg/ai"
	"github.com/redhat-et/copilot-ops/pkg/ai/bloom"
	"github.com/redhat-et/copilot-ops/pkg/ai/gpt3"
	"github.com/redhat-et/copilot-ops/pkg/ai/gptj"
	"github.com/redhat-et/copilot-ops/pkg/cmd/config"
	"github.com/redhat-et/copilot-ops/pkg/filemap"
	"github.com/spf13/cobra"
)

// Request Defines the necessary values used when requesting new files from the selected
// AI backends.
// FIXME: consolidate the settings depending on the type of Model. E.g., OpenAI settings should be under their own.
type Request struct {
	Config       config.Config
	Fileset      *config.Filesets
	Filemap      *filemap.Filemap
	FilemapText  string
	UserRequest  string
	IsWrite      bool
	OutputType   string
	OpenAIURL    string
	NTokens      int32
	NCompletions int32
	// Backend Sepecifies which type of AI Backend to use.
	Backend ai.Backend
}

// PrepareRequest Processes the user input along with provided environment variables,
// creating a Request object which is used for context in further requests.
func PrepareRequest(cmd *cobra.Command) (*Request, error) {
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
	aiBackend, _ := cmd.Flags().GetString(FlagAIBackendFull)

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
	conf := config.Config{}
	if err := conf.Load(); err != nil {
		return nil, err
	}

	// load files
	fm := filemap.NewFilemap()
	if err := fm.LoadFiles(files); err != nil {
		log.Fatalf("error loading files: %s\n", err.Error())
	}
	if len(filesets) > 0 {
		log.Printf("loading filesets: %v\n", filesets)
	}
	if err := fm.LoadFilesets(filesets, conf, config.ConfigFile); err != nil {
		log.Fatalf("error loading filesets: %s\n", err.Error())
	}
	filemapText := fm.EncodeToInputText()

	// select backend type
	selectedBackend := ai.Backend(aiBackend)
	if selectedBackend == "" {
		selectedBackend = conf.Backend
	}

	// configure OpenAI
	// FIXME: create default config methods for these
	if openAIURL != "" {
		conf.OpenAI.BaseURL = openAIURL
	}

	// configure GPT-J
	if conf.GPTJ == nil {
		conf.GPTJ = &gptj.Config{
			URL: gptj.APIURL,
		}
	}

	// default for bloom
	if conf.BLOOM == nil {
		conf.BLOOM = &bloom.Config{
			URL: bloom.APIURL,
		}
	}

	r := Request{
		Config:       conf,
		Filemap:      fm,
		FilemapText:  filemapText,
		UserRequest:  request,
		IsWrite:      write,
		OutputType:   outputType,
		NTokens:      nTokens,
		NCompletions: nCompletions,
		Backend:      selectedBackend,
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

	cmd.Flags().StringP(
		FlagAIBackendFull, FlagAIBackendShort, string(ai.GPT3), "AI Backend to use",
	)

	_ = cmd.Flags().StringP(
		FlagOpenAIURLFull,
		FlagOpenAIURLShort,
		gpt3.OpenAIURL+gpt3.OpenAIEndpointV1,
		"OpenAI URL",
	)
}
