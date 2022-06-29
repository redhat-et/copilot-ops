// types.go Defines the types within the openai package.
package openai

import "net/http"

// EditResponse Represents the response body from OpenAI for requests to /edit.
type EditResponse struct {
	OpenAIResponse
}

type OpenAIResponse struct {
	Object  string `json:"object"`
	Created uint64 `json:"created"`
	Choices []struct {
		Text  string `json:"text"`
		Index int    `json:"index"`
	} `json:"choices"`
}

type BodyParameters struct {
	// Model is the model used by OpenAI to generate completions
	Model string `json:"model"`
	// Temperature Is the sampling temperature to use. Higher values means the model will take more risks. Try 0.9 for more creative applications, and 0 (argmax sampling) for ones with a well-defined answer.
	Temperature float32 `json:"temperature,omitempty"`
	TopP        int     `json:"top_p,omitempty"`
}

// EditRequestBody Defines the parameters for the /edit endpoint.
type EditRequestBody struct {
	BodyParameters
	Instruction string `json:"instruction"`
	Input       string `json:"input,omitempty"`
}

// CompletionRequestBody Defines the body of data for requests that must be sent to the Completions endpoint.
type CompletionRequestBody struct {
	BodyParameters
	// Prompt is the prompt to generate completions for, it can be encoded as a string, or as an array of strings.
	Prompt string `json:"prompt"`
	// Suffix is the suffix that comes after a completion of inserted text.
	Suffix string `json:"suffix,omitempty"`
	// MaxTokens Is an integer representing the maximum number of tokens to generate in the completion.
	MaxTokens int32 `json:"max_tokens,omitempty"`
	// N Is a number indicating how many completions should be generated for the given prompt.
	N int32 `json:"n,omitempty"`
	// Logprobs Is an integer which determines how many log probabilities OpenAI will return for its most likely tokens.
	// Maximum is 5.
	Logprobs int32 `json:"logprobs,omitempty"`
	// Echo is a boolean indicating whether we want to echo the prompt in addition to the completion.
	Echo bool `json:"echo,omitempty"`
	// Stop Is a list of up to 4 sequences where the API will stop generating further tokens.
	Stop []string `json:"stop,omitempty"`
	// PresencePenalty Number between -2.0 and 2.0. Positive values penalize new tokens based on whether they appear in the text so far, increasing the model's likelihood to talk about new topics.
	PresencePenalty float32 `json:"presence_penalty,omitempty"`
	// FrequencyPenalty Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing frequency in the text so far, decreasing the model's likelihood to repeat the same line verbatim.
	FrequencyPenalty float32 `json:"frequency_penalty,omitempty"`
	// BestOf Generates best_of completions server-side and returns the "best" (the one with the lowest log probability per token). Results cannot be streamed.
	BestOf int `json:"best_of,omitempty"`
	// User Is a string which uniquely identifies a user performing the completions.
	User string `json:"user,omitempty"`
}

// CompletionResponse Is the object returned from OpenAI after a successful request has been made to /completions
type CompletionResponse struct {
	OpenAIResponse
	// ID Is a string which uniquely identifies this completion for OpenAI.
	ID string `json:"id,omitempty"`
	// Model Identifies the model that was used in order to make this completion.
	Model string `json:"model,omitempty"`
}

// SearchRequestBody Defines the parameters for the /search endpoint.
type SearchRequestBody struct {
	// Query Is the Query to search against documents.
	Query string `json:"query"`
	// Documents Is an array of documents to search against.
	// If this is specified then File must be omitted.
	Documents []string `json:"documents,omitempty"`
	// File Is the path to a file containing documents to search against.
	// If this is specified then Documents must be omitted.
	File string `json:"file,omitempty"`
	// MaxRerank Is the maximum number of documents to be re-ranked and returned from search.
	MaxRerank int32 `json:"max_rerank,omitempty"`
	// ReturnMetadata A special boolean flag for showing metadata.
	// If set to true, each document entry in the returned JSON will contain a "metadata" field.
	// This flag only takes effect when file is set.
	ReturnMetadata bool `json:"return_metadata,omitempty"`
	// User Is a unique identifier representing your end-user, which will help OpenAI to monitor and detect abuse.
	User string `json:"user,omitempty"`
}

// SearchResponseBody Represents the response body from OpenAI for requests to /search.
type SearchResponseBody struct {
	// Data Returns a ranking of the documents in the search results.
	Data []struct {
		// Document Is the document that was searched against.
		Document int32 `json:"document"`
		// Score Is the score of the document.
		Score float64 `json:"score"`
		// Object Defines the type of object that is being returned.
		Object string `json:"object"`
	} `json:"data"`
	// Object Defines the type of object that is being returned.
	Object string `json:"object"`
}

// OpenAIClient Defines the client for interacting with the OpenAI API
type OpenAIClient struct {
	// Client is the HTTP client used to perform requests to the OpenAI API.
	Client *http.Client
	// APIUrl Defines the endpoint that the client will use to interact with the OpenAI API.
	APIUrl string
	// Engine represents the engine used when getting completions from OpenAI.
	// TODO: Create a constant enumeration of engines so that we can include their character limits.
	Engine string
	// AuthToken is the token used to authenticate with the OpenAI API.
	AuthToken string
	// OrganizationID is the ID of the organization used to authenticate with the OpenAI API.
	OrganizationID string
	// NTokens Is an integer representing the maximum number of tokens to generate in the completion.
	NTokens int32
	// NCompletions is an integer representing the maximum number of completions to generate from a prompt.
	NCompletions int32
}

// ErrorResponse Defines an Error object which will be returned by OpenAI when an error occurs.
type ErrorResponse struct {
	Error *struct {
		Code    *int    `json:"code,omitempty"`
		Message string  `json:"message"`
		Type    *string `json:"type,omitempty"`
		Param   string  `json:"param"`
	} `json:"error,omitempty"`
}

// FileOutput Defines a file object used to store generated file attributes.
type FileOutput struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Contents string `json:"contents"`
}

// GeneratedFilesOutput Defines a list of all generated files.
type GeneratedFilesOutput struct {
	GeneratedFiles []FileOutput
}
