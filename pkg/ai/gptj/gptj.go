package gptj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redhat-et/copilot-ops/pkg/ai"
)

// Define constants used by EleutherAI here.
const (
	APIURL             = "https://api.eleuther.ai"
	CompletionEndpoint = "completion"
)

// GenerateParams Defines the parameters which are sent when requesting
// a response from Eleuther's publicly hosted GPT-J instance on https://6b.eleuther.ai.
//   --data-raw '{"context":"asdfasdf","top_p":0.9,"temp":0.8,"response_length":128,"remove_input":true}' \
type GenerateParams struct {
	// Context Defines the prompt which will be passed into GPT-J.
	Context        string  `json:"context"`
	TopP           float32 `json:"top_p"`
	Temp           float32 `json:"temp"`
	ResponseLength int32   `json:"response_length"`
	RemoveInput    bool    `json:"remove_input"`
}

// Generate Invokes the generate function to GPT-J. Currently, the endpoint
// only supports a single item to be returned when generated.
func (c gptjClient) Generate(params interface{}) ([]string, error) {
	// parse params
	generateParams, ok := params.(GenerateParams)
	if !ok {
		return nil, fmt.Errorf("could not parse params")
	}

	// marshal params into json bytes
	reqBytes, err := json.Marshal(generateParams)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
	reqBuff := bytes.NewBuffer(reqBytes)

	// create request
	urlPath := c.baseUrl + "/" + CompletionEndpoint
	req, err := http.NewRequest(http.MethodPost, urlPath, reqBuff)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	// transform request into response
	var response []choice
	if err = c.sendRequest(req, &response); err != nil {
		return nil, fmt.Errorf("could not request eleuther.ai: %w", err)
	}
	choices := make([]string, len(response))
	for i, choice := range response {
		choices[i] = choice.GeneratedText
	}
	return choices, nil
}

func (c gptjClient) Edit(params interface{}) ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateGPTJClient Returns a GPT-J client which implements the AI Client interface.
func CreateGPTJClient(url string, client *http.Client) ai.Client {
	c := gptjClient{
		baseUrl:    url,
		httpClient: client,
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	return c
}
