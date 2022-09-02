// ai declares an interface for accessing various AI backends.
package ai

// Client Is an interface which declares common methods expected of
// various NLP models; regardless of the architecture style or method
// of accessing.
type Client interface {
	// Generate Returns a list of string which are potential completions
	// for the given prompt.
	Generate(params interface{}) ([]string, error)
	// Edit Returns a list of edits made to the provided data.
	Edit(params interface{}) ([]string, error)
}
