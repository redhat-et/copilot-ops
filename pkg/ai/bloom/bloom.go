package bloom

import (
	"fmt"

	"github.com/redhat-et/copilot-ops/pkg/ai"
)

type bloomClient struct {
}

func (c bloomClient) Generate() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c bloomClient) Edit() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateOPTGenerateClient Returns an OPT-175B client capable of making code generations.
func CreateOPTGenerateClient() ai.GenerateClient {
	return bloomClient{}
}

// CreateOPTEditClient Returns an OPT-175B client capable of making code edits.
func CreateOPTEditClient() ai.EditClient {
	return bloomClient{}
}
