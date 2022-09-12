// ai declares an interface for accessing various AI backends.
package ai

// Client Is an interface which declares common methods expected of
// various NLP models; regardless of the architecture style or method
// of accessing.
type Client interface {
	Generate(params interface{}) ([]string, error)
	Edit(params interface{}) ([]string, error)
}

// GenerateClient Describes a client which can generate code/text from an AI backend.
type GenerateClient interface {
	// Generate Returns a list of string which are potential completions
	// for the given prompt.
	Generate() ([]string, error)
}

// EditClient Describes an AI client capable of implementing the edit function.
type EditClient interface {
	// Edit Returns a list of edits made to the provided data.
	Edit() ([]string, error)
}

// Backend Defines a type specifically for backends.
type Backend string

const (
	// GPT3 Declares the GPT-3 AI backend, created and hosted by OpenAI.
	GPT3 Backend = "gpt-3"
	// GPTJ Declares the GPT-J AI backend, created and hosted by EleutherAI.
	GPTJ Backend = "gpt-j"
	// BLOOM Declares the BLOOM AI backend, created by BigScience, hosted by HuggingFace.
	BLOOM Backend = "bloom"
	// OPT Declares the OPT-175B AI Backend, created by Meta.
	OPT Backend = "opt"
	// Unselected Represents an empty AI backend type.
	Unselected Backend = ""
)
