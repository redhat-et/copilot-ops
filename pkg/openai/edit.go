// Package openai defines functions for interacting with OpenAI's endpoints
package openai

const (
	EditEndpoint string = "https://api.openai.com/v1/engines/code-davinci-002/edits"
)

// EditCode accepts an input and an instruction string and returns an
// output string edited by the AI engine.
func EditCode(input string, instruction string) (string, error) {

	// TODO for now - return the input unedited
	return input, nil

}
