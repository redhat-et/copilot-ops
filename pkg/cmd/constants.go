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
	FlagOutputTypeFull    = "output"
	FlagOutputTypeShort   = "o"
)

// COMMAND Constants which define the names of commands used in the CLI.
const (
	CommandEdit     = "edit"
	CommandGenerate = "generate"
)

// Miscellaneous constants used in the CLI.
const (
	DefaultTokens      = 512
	DefaultCompletions = 1
)
