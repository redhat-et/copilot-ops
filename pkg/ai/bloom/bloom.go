package bloom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redhat-et/copilot-ops/pkg/ai"
)

const (
	// APIURL Defines where you can find the BLOOM API.
	APIURL = "https://api-inference.huggingface.co/models/bigscience/bloom"
	// DefaultTokenSize Defines the default amount of max tokens set by BLOOM.
	DefaultTokenSize = 100
)

// Config Defines the values required for successful connections to BLOOM 176B.
type Config struct {
	// URL Defines where to find the API.
	URL string `json:"url" yaml:"url"`
}

// GenerateParameters Defines a struct which sets the parameters to BLOOM.
type GenerateParameters struct {
	Seed          int     `json:"seed"`
	EarlyStopping bool    `json:"early_stopping"`
	LengthPenalty int     `json:"length_penalty"`
	MaxNewTokens  int     `json:"max_new_tokens"`
	DoSample      bool    `json:"do_sample"`
	TopP          float32 `json:"top_p"`
}

// bloomClient Describes the client which implements the AI interfaces.
type bloomClient struct {
	BaseURL        string
	generateParams *generateRequest
}

// generateRequest Defines the body of a request to BLOOM's completions endpoint.
type generateRequest struct {
	Inputs     string             `json:"inputs"`
	Parameters GenerateParameters `json:"parameters"`
}

// choice Represents a single text-generation element.
type choice struct {
	GeneratedText string `json:"generated_text"`
}

// responseError Represents an object returned in the event of an error.
type responseError struct {
	Error *string `json:"error,omitempty"`
}

// Generate Returns a list of completions created by the BLOOM BigModel API.
func (c bloomClient) Generate() ([]string, error) {
	if c.generateParams == nil {
		return nil, fmt.Errorf("no params provided")
	}

	// marshal params into json bytes
	reqBytes, err := json.Marshal(c.generateParams)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
	reqBuff := bytes.NewBuffer(reqBytes)
	if reqBuff == nil {
		return nil, fmt.Errorf("could not initialize buffer")
	}

	// create request
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, c.BaseURL, reqBuff)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("content-type", "application/json")

	// retrieve a response from server
	resp := make([]choice, 0)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	// wrap the HTTP error
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("error, status code: %d", res.StatusCode)
	}

	// attempt to marshal into a response
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		// attempt to unmarshal error
		errResp := responseError{}
		err2 := json.NewDecoder(res.Body).Decode(&errResp)
		if err2 != nil || errResp.Error == nil {
			return nil, fmt.Errorf("could not read response: %w", err)
		}

		return nil, fmt.Errorf("received error from server: %s", *errResp.Error)
	}

	// transform into flat list
	responses := make([]string, len(resp))
	for i, gen := range resp {
		responses[i] = gen.GeneratedText
	}
	return responses, nil
}

// Edit Returns a list of edits made by the BLOOM BigModel API.
func (c bloomClient) Edit() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateBloomGenerateClient Returns a client which represents a request made to the OpenAI API.
func CreateBloomGenerateClient(conf Config, prompt string, params GenerateParameters) ai.GenerateClient {
	genParams := &generateRequest{
		Inputs:     prompt,
		Parameters: params,
	}
	return bloomClient{
		BaseURL:        conf.URL,
		generateParams: genParams,
	}
}

// CreateBloomEditClient Returns a client capable of making edits to the OpenAI API.
func CreateBloomEditClient() ai.EditClient {
	return bloomClient{}
}
