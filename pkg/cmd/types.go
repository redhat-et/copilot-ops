// Package cmd defines the functions for the CLI.
package cmd

import (
	fm "github.com/redhat-et/copilot-ops/pkg/filemap"
)

// Result represents the command results, warnings, and errors if
// they exist. Warnings will not block a result from being returned, but errors will.
type Result struct {
	Error    *Error        `json:"error"`
	Warnings *[]Error      `json:"warnings"`
	Result   *[]fm.Filemap `json:"result"`
}

// Error represents an error or warning from within the program
// which has taken place during the execution of the CLI.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    error  `json:"data"`
}

type Warning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    error  `json:"data"`
}
