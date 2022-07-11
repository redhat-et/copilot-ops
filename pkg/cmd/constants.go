package cmd

// Define the names of flags used in commands.
const (
	FlagRequest      = "request"
	FlagWrite        = "write"
	FlagPath         = "path"
	FlagFiles        = "file"
	FlagFilesets     = "fileset"
	FlagNTokens      = "ntokens"
	FlagNCompletions = "ncompletions"
	FlagOutputType   = "output"
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
