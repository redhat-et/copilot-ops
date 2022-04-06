// Package openai defines functions for interacting with OpenAI's endpoints
package openai

import (
	fm "github.com/redhat-et/copilot-ops/pkg/filemap"
)

const (
	EditEndpoint string = "https://api.openai.com/v1/engines/code-davinci-002/edits"
)

// GetFileEditsFromOpenAI Accepts a Filemap and an instruction string to edit the provided files
// using a response obtained from OpenAI.
func GetFileEditsFromOpenAI(fileMap fm.Filemap, instruction string) (EditReply, error) {
	return EditReply{}, nil
}

// ResponseToFilemap Parses OpenAI's response to populate the Filemap
func ResponseToFilemap(response EditReply) (fm.Filemap, error) {
	return fm.Filemap{}, nil
}
