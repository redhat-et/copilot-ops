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
	// TODO: Create a constant enumeration of engines so that we can include their character limits.
	Engine string
	// AuthToken is the token used to authenticate with the OpenAI API.
	AuthToken string
	// OrganizationID is the ID of the organization used to authenticate with the OpenAI API.
	OrganizationID string
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
