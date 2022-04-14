// types.go Defines the types within the openai package.
package openai

import "net/http"

// EditResponseBody Represents the response body from OpenAI for requests to /edit.
type EditResponseBody struct {
	Object  string `json:"object"`
	Created uint64 `json:"created"`
	Choices []struct {
		Text  string `json:"text"`
		Index int    `json:"index"`
	} `json:"choices"`
}

// EditRequestBody Defines the parameters for the /edit endpoint.
type EditRequestBody struct {
	Instruction string  `json:"instruction"`
	Input       string  `json:"input,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	TopP        int     `json:"top_p,omitempty"`
}

// OpenAIClient Defines the client for interacting with the OpenAI API
type OpenAIClient struct {
	// Client is the HTTP client used to perform requests to the OpenAI API.
	Client *http.Client
	// APIUrl Defines the endpoint that the client will use to interact with the OpenAI API.
	APIUrl string
	// Engine represents the engine used when getting completions from OpenAI.
	Engine string
}
