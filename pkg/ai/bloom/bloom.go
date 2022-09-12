package bloom

import (
	"fmt"
	"net/http"

	"github.com/redhat-et/copilot-ops/pkg/ai"
)

// Config Defines the values required for successful connections to BLOOM 176B.
type Config struct {
	// URL Defines where to find the API.
	URL string
	// HTTPClient is the client which will be used when making requests.
	HTTPClient *http.Client
}

type bloomClient struct {
}

// Generate Returns a list of completions created by the BLOOM BigModel API.
func (c bloomClient) Generate() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

// Edit Returns a list of edits made by the BLOOM BigModel API.
func (c bloomClient) Edit() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateBloomGenerateClient Returns a client which represents a request made to the OpenAI API.
func CreateBloomGenerateClient() ai.GenerateClient {
	return bloomClient{}
}

// CreateBloomEditClient Returns a client capable of making edits to the OpenAI API.
func CreateBloomEditClient() ai.EditClient {
	return bloomClient{}
}
