// Package openai defines functions for interacting with OpenAI's endpoints
package openai

import (
	fm "github.com/redhat-et/copilot-ops/pkg/filemap"
)

// EditRequest represents user-information which will be passed to OpenAI's edit endpoint.
type EditRequest struct {
	Filemap     *fm.Filemap `json:"filemap"`
	Instruction string      `json:"instruction"`
}

// EditReply represents the response from OpenAI's edit endpoint.
type EditReply struct {
	Filemap *fm.Filemap `json:"filemap"`
}

// GetFileEditsFromOpenAI Accepts a Filemap and an instruction string to edit the provided files
// using a response obtained from OpenAI.
func GetFileEditsFromOpenAI(fileMap fm.Filemap, instruction string) (EditReply, error) {
	return EditReply{}, nil
}

// ResponseToFilemap Parses OpenAI's response to populate the Filemap
func ResponseToFilemap(response EditReply) (fm.Filemap, error) {
	return fm.Filemap{}, nil
}
