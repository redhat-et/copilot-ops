// types.go Defines the types within the openai package.
package openai

import fm "github.com/redhat-et/copilot-ops/pkg/filemap"

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

// EditRequest represents user-information which will be passed to OpenAI's edit endpoint.
type EditRequest struct {
	Filemap     *fm.Filemap `json:"filemap"`
	Instruction string      `json:"instruction"`
}

// EditReply represents the response from OpenAI's edit endpoint.
type EditReply struct {
	Filemap *fm.Filemap `json:"filemap"`
}
