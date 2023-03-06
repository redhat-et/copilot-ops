package gpt3

import (
	"context"
	"fmt"

	"github.com/redhat-et/copilot-ops/pkg/ai"
	gogpt "github.com/sashabaranov/go-openai"
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

// gpt3Client Is a wrapper struct around the go-openai
// package.
type gpt3Client struct {
	client           gogpt.Client
	editParams       *gogpt.EditsRequest
	completionParams *gogpt.CompletionRequest
}

// Config Defines the values required for connecting to the GPT-3 API.
// FIXME: set better names for these fields.
type Config struct {
	// APIKey Is the API token used when making requests.
	APIKey string `json:"apiKey" yaml:"apiKey"`
	// OrgID Is an optional value which is set by users to dictate billing information.
	OrgID *string `json:"orgID,omitempty" yaml:"orgID,omitempty"`
	// BaseURL Defines where the client will reach out to contact the API.
	BaseURL string `json:"url" yaml:"url"`
}

// Generate Reaches out to the OpenAI GPT-3 Completions API and returns
// a list of completions pertinent to the request.
func (c gpt3Client) Generate() ([]string, error) {
	if c.completionParams == nil {
		return nil, fmt.Errorf("no completions params were provided")
	}
	// make request
	resp, err := c.client.CreateCompletion(context.TODO(), *c.completionParams)
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
func (c gpt3Client) Edit() ([]string, error) {
	// ensure params
	if c.editParams == nil {
		return nil, fmt.Errorf("no edit params were provided")
	}
	// editParams, ok := params.(EditParams)
	resp, err := c.client.Edits(context.Background(), *c.editParams)
	if err != nil {
		return nil, fmt.Errorf("could not request openai: %w", err)
	}
	edits := make([]string, len(resp.Choices))
	for i, choice := range resp.Choices {
		edits[i] = choice.Text
	}
	return edits, nil
}

// CreateGPT3GenerateClient Returns a GPT-3 client which accesses OpenAI's
// GPT-3 endpoint to generate completions.
func CreateGPT3GenerateClient(conf Config, prompt string, maxTokens, nCompletions int) ai.GenerateClient {
	// create a GPT-3 Client
	client := createGPT3Client(conf)
	// create params for getting a completion
	params := &gogpt.CompletionRequest{
		Model:       OpenAICodeDavinciV2,
		Prompt:      prompt,
		MaxTokens:   maxTokens,
		N:           nCompletions,
		Temperature: 0.0,
		Stop:        []string{CompletionEndOfSequence},
	}

	return gpt3Client{
		client:           *client,
		completionParams: params,
	}
}

// CreateGPT3EditClient Returns a client based on GPT-3 capable of performing edits.
func CreateGPT3EditClient(
	conf Config,
	input, instruction string,
	numEdits int, temperature,
	topP *float32,
) ai.EditClient {
	client := createGPT3Client(conf)
	// set params
	model := OpenAICodeDavinciEditV1
	editParams := &gogpt.EditsRequest{
		Model:       &model,
		N:           numEdits,
		Input:       input,
		Instruction: instruction,
	}
	if temperature != nil {
		editParams.Temperature = *temperature
	}
	if topP != nil {
		editParams.TopP = *topP
	}

	return gpt3Client{
		client:     *client,
		editParams: editParams,
	}
}

// createGPT3Client Returns a go-gpt client using the provided config.
func createGPT3Client(conf Config) *gogpt.Client {
	orgID := ""
	if conf.OrgID != nil {
		orgID = *conf.OrgID
	}
	return gogpt.NewClientWithConfig(gogpt.ClientConfig{
		BaseURL: conf.BaseURL,
		OrgID:   orgID,
	})
}
