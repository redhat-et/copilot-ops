// Package cmd defines the functions for the CLI.
package cmd

import fm "github.com/redhat-et/copilot-ops/pkg/filemap"

// CmdResult represents the command results, warnings, and errors if
// they exist. Warnings will not block a result from being returned, but errors will.
type CmdResult struct {
	Errors   *[]CmdError   `json:"errors"`
	Warnings *[]CmdError   `json:"warnings"`
	Result   *[]fm.Filemap `json:"result"`
}

// CmdError represents an error or warning from within the program
// which has taken place during the execution of the CLI.
// CmdError is based on the JSON-RPC error type,
// 	see: https://www.jsonrpc.org/specification#error_object.
type CmdError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    error  `json:"data"`
}
