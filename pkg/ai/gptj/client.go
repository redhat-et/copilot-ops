package gptj

import (
	"net/http"
)

// gptjClient is a client implementation of GPT-J meant to implement
// the AI Client interface.
type gptjClient struct {
	baseURL        string
	httpClient     *http.Client
	generateParams *GenerateParams
}

// choice Represents a list element which is returned from the GPT-J API.
type choice struct {
	GeneratedText string `json:"generated_text"`
}
