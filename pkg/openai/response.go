// Package openai defines functions for interacting with OpenAI's endpoints
package openai

import (
	"fmt"
)

// GetFirstChoice Returns the string contents of the first choice returned by the AI engine.
func (resp *OpenAIResponse) GetFirstChoice() (string, error) {
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Text, nil
	}
	return "", fmt.Errorf("no choices found")
}

// GetAllChoices returns a list containing all of the choices returned by the AI engine.
func (resp *OpenAIResponse) GetAllChoices() []string {
	var choices []string
	for _, choice := range resp.Choices {
		choices = append(choices, choice.Text)
	}
	return choices
}

// GetNChoices returns a list of n choices returned by the AI engine.
// n is the number of completions requested by the user.
func (resp *OpenAIResponse) GetNChoices(n int32) ([]string, error) {
	var choices []string
	var i int32

	if len(resp.Choices) == 0 {
		return choices, fmt.Errorf("no choices found")
	}
	for i = 0; i < n; i++ {
		choices = append(choices, resp.Choices[i].Text)
	}

	return choices, nil
}
