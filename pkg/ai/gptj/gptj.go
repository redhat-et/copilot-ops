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
type GenerateParams struct {
	// Context Defines the prompt which will be passed into GPT-J.
	Context        string  `json:"context"`
	TopP           float32 `json:"top_p"`
	Temp           float32 `json:"temp"`
	ResponseLength int32   `json:"response_length"`
	RemoveInput    bool    `json:"remove_input"`
}

// Config Describes the structure needed for configuring a GPT-J client which
// connects to Eleuther AI's endpoints.
type Config struct {
	// URL Defines the URL which the HTTP Client will be making requests to.
	URL string
	// HTTPClient Is an HTTP Client which is used in making requests.
	HTTPClient *http.Client
}

// Generate Invokes the generate function to GPT-J. Currently, the endpoint
// only supports a single item to be returned when generated.
func (c gptjClient) Generate() ([]string, error) {
	// parse params
	if c.generateParams == nil {
		return nil, fmt.Errorf("no params provided")
	}

	// marshal params into json bytes
	reqBytes, err := json.Marshal(c.generateParams)
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

// CreateGPTJGenerateClient Returns a GPT-J client which implements the AI Client interface.
func CreateGPTJGenerateClient(conf Config, params GenerateParams) ai.GenerateClient {
	httpClient := conf.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := gptjClient{
		baseUrl:    conf.URL,
		httpClient: httpClient,
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	return c
}
