// Package openai defines functions for interacting with OpenAI's endpoints
package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	EditEndpoint string = "https://api.openai.com/v1/engines/code-davinci-002/edits"
)

// EditCode accepts an input and an instruction string and returns an
// output string edited by the AI engine.
func EditCode(input string, instruction string) (string, error) {
	// build the edit body
	var editBody EditRequestBody = EditRequestBody{
		Instruction: instruction,
		Input:       input,
	}

	// marshal the edit body into a JSON string
	jsonEditBody, err := json.Marshal(editBody)
	if err != nil {
		return "", err
	}

	// build the edit request
	body := io.Reader(bytes.NewReader(jsonEditBody))

	req, err := http.NewRequest("POST", EditEndpoint, body)
	if err != nil {
		return "", err
	}

	// set the request headers
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "openai-cli")

	// build the client
	client := &http.Client{
		// NOTE: This could vary depending on the length of input & options used
		Timeout: time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("edit request failed with status %d", resp.StatusCode)
	}

	// decode the response body
	var editResponse EditResponseBody
	err = json.NewDecoder(resp.Body).Decode(&editResponse)
	if err != nil {
		return "", err
	}

	// get the first response
	response, err := editResponse.GetFirstResponse()
	if err != nil {
		return "", err
	}

	// TODO for now - return the input unedited
	return response, nil
}

func (erb *EditResponseBody) GetFirstResponse() (string, error) {
	if len(erb.Choices) > 0 {
		return erb.Choices[0].Text, nil
	}
	return "", fmt.Errorf("no response found")
}
