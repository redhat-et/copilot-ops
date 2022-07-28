package openai

const (
	EditEndpoint       string = "edits"
	CompletionEndpoint string = "completions"
	SearchEndpoint     string = "search"
	OpenAIURL          string = "https://api.openai.com"
	// Maybe the OpenAIEndpoint should be a part of the URL string?
	OpenAIEndpointV1        string = "/v1"
	OpenAICodeDavinciEditV1 string = "code-davinci-edit-001"
	OpenAICodeDavinciV2     string = "code-davinci-002"
	CompletionEndOfSequence string = "EOF"
)
