package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// CreateOpenAIClient Creates a client to perform requests to OpenAI.
func CreateOpenAIClient() (*OpenAIClient, error) {
	return &OpenAIClient{
		Client: &http.Client{
			Timeout: time.Minute,
		},
		APIUrl: OpenAIEndpointV1,
		Engine: OpenAICodeDavinci2,
	}, nil
}

// GetAPIUrl Creates a URL for the OpenAI API.
func (openAI OpenAIClient) EnginePath() string {
	return fmt.Sprintf("engines/%s", openAI.Engine)
}

// APIHeaders Returns a map of headers to be used when making requests to the OpenAI API.
func (openAI *OpenAIClient) APIHeaders() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + os.Getenv("OPENAI_API_KEY"),
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"User-Agent":    "copilot-ops-cli",
	}
}

// MakeRequest Makes a request to the OpenAI API, suffixed with the given endpoint with the provided
// body, returning the response as a stream of bytes to be unmarshaled by the caller. An error is
// returned if the request fails.
func (openAI *OpenAIClient) MakeRequest(endpoint string, body interface{}) ([]byte, error) {
	// marshal the struct into json
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// stringify the body as bytes
	bodyBytes := io.Reader(bytes.NewReader(jsonBody))

	// build the request
	apiURL := openAI.APIUrl + "/" + openAI.EnginePath() + "/" + endpoint
	req, err := http.NewRequest("POST", apiURL, bodyBytes)
	if err != nil {
		return nil, err
	}

	// set the request headers
	for k, v := range openAI.APIHeaders() {
		req.Header.Set(k, v)
	}

	// make the request
	resp, err := openAI.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

// EditCode accepts an input and an instruction string and returns an
// output string edited by the AI engine.
func (openAI *OpenAIClient) EditCode(input string, instruction string) (string, error) {
	// build the edit body
	var editBody EditRequestBody = EditRequestBody{
		Instruction: instruction,
		Input:       input,
	}

	// make the request
	data, err := openAI.MakeRequest(EditEndpoint, editBody)
	if err != nil {
		return "", err
	}

	// decode the response body
	var editResponse EditResponseBody
	err = json.Unmarshal(data, &editResponse)
	if err != nil {
		return "", err
	}

	// get the first response
	response, err := editResponse.GetFirstChoice()
	if err != nil {
		return "", err
	}

	// TODO for now - return the input unedited
	return response, nil
}
