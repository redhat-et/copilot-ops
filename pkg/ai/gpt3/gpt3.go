package gpt3

import (
	"context"
	"fmt"

	"github.com/redhat-et/copilot-ops/pkg/ai"
	gogpt "github.com/sashabaranov/go-gpt3"
)

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

type GenerateParams struct {
	Params gogpt.CompletionRequest
}

type EditParams struct {
	Params gogpt.EditsRequest
}

// gpt3Client Is a wrapper struct around the go-gpt3
// package.
type gpt3Client struct {
	client gogpt.Client
}

// Generate Reaches out to the OpenAI GPT-3 Completions API and returns
// a list of completions pertinent to the request.
func (c gpt3Client) Generate(params interface{}) ([]string, error) {
	genParams, ok := params.(GenerateParams)
	if !ok {
		return nil, fmt.Errorf("could not parse params")
	}
	// make request
	resp, err := c.client.CreateCompletion(context.Background(), genParams.Params)
	if err != nil {
		return nil, err
	}
	// collect strings from response
	responses := make([]string, len(resp.Choices))
	for i, choice := range resp.Choices {
		responses[i] = choice.Text
	}
	return responses, nil
}

// Edit Reaches out to the OpenAI GPT-3 Edits API and returns a list of
// responses which have been edited in accordance with the given instruction.
func (c gpt3Client) Edit(params interface{}) ([]string, error) {
	editParams, ok := params.(EditParams)
	if !ok {
		return nil, fmt.Errorf("could not parse params")
	}
	resp, err := c.client.Edits(context.Background(), editParams.Params)
	if err != nil {
		return nil, fmt.Errorf("could not request openai: %w", err)
	}
	edits := make([]string, len(resp.Choices))
	for i, choice := range resp.Choices {
		edits[i] = choice.Text
	}
	return edits, nil
}

// CreateGPT3Client Returns a struct based on OpenAI which implements the
// AIModel interface.
func CreateGPT3Client(token string, orgID *string, url string) ai.Client {
	var client *gogpt.Client
	if orgID != nil {
		client = gogpt.NewOrgClient(token, *orgID)
	} else {
		client = gogpt.NewClient(token)
	}
	client.BaseURL = url
	return gpt3Client{
		client: *client,
	}
}
