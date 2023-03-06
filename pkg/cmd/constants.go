package cmd

// Define the names of flags used in commands.
const (
	FlagRequestFull       = "request"
	FlagRequestShort      = "r"
	FlagWriteFull         = "write"
	FlagWriteShort        = "w"
	FlagPathFull          = "path"
	FlagPathShort         = "p"
	FlagFilesFull         = "file"
	FlagFilesShort        = "f"
	FlagFilesetsFull      = "fileset"
	FlagFilesetsShort     = "s"
	FlagNTokensFull       = "ntokens"
	FlagNTokensShort      = "n"
	FlagNCompletionsFull  = "ncompletions"
	FlagNCompletionsShort = "c"
	FlagOpenAIURLFull     = "openai-url"
	FlagOpenAIURLShort    = "d"
	FlagOutputTypeFull    = "output"
	FlagOutputTypeShort   = "o"
	FlagAIBackendFull     = "backend"
	FlagAIBackendShort    = "b"
)

// COMMAND Constants which define the names of commands used in the CLI.
const (
	CommandEdit     = "edit"
	CommandGenerate = "generate"
	CommandAsk      = "ask"
)

// Miscellaneous constants used in the CLI.
const (
	DefaultTokens      = 1000
	DefaultCompletions = 1
)
