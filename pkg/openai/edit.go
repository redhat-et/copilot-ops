// Package openai defines functions for interacting with OpenAI's endpoints
package openai

import (
	"fmt"
)

const (
	EditEndpoint            string = "edits"
	CompletionEndpoint      string = "completions"
	OpenAIEndpointV1        string = "https://api.openai.com/v1"
	OpenAICodeDavinciEditV1 string = "code-davinci-edit-001"
)

// GetFirstChoice Returns the string contents of the first choice returned by the AI engine.
func (resp *EditResponseBody) GetFirstChoice() (string, error) {
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Text, nil
	}
	return "", fmt.Errorf("no choices found")
}

// GetAllChoices returns a list containing all of the choices returned by the AI engine.
func (resp *EditResponseBody) GetAllChoices() []string {
	var choices []string
	for _, choice := range resp.Choices {
		choices = append(choices, choice.Text)
	}
	return choices
}
