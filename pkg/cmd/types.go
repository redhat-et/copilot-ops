// Package cmd defines the functions for the CLI.
package cmd

import (
	fm "github.com/redhat-et/copilot-ops/pkg/filemap"
)

// CmdResult represents the command results, warnings, and errors if
// they exist. Warnings will not block a result from being returned, but errors will.
type CmdResult struct {
	Error    *CmdError     `json:"error"`
	Warnings *[]CmdError   `json:"warnings"`
	Result   *[]fm.Filemap `json:"result"`
}

// CmdError represents an error or warning from within the program
// which has taken place during the execution of the CLI.
type CmdError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    error  `json:"data"`
}

type CmdWarning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    error  `json:"data"`
}
