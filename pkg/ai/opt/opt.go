package opt

import (
	"fmt"

	"github.com/redhat-et/copilot-ops/pkg/ai"
)

type optClient struct {
}

func (c optClient) Generate() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c optClient) Edit() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateOPTGenerateClient Returns an OPT-175B client capable of making code generations.
func CreateOPTGenerateClient() ai.GenerateClient {
	return optClient{}
}

// CreateOPTEditClient Returns an OPT-175B client capable of making code edits.
func CreateOPTEditClient() ai.EditClient {
	return optClient{}
}
