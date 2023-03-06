package gptj

import (
	"fmt"

	"github.com/redhat-et/copilot-ops/pkg/ai"
)

type gptjClient struct {
}

func (c gptjClient) Generate() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c gptjClient) Edit() ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

// CreateOPTGenerateClient Returns an OPT-175B client capable of making code generations.
func CreateOPTGenerateClient() ai.GenerateClient {
	return gptjClient{}
}

// CreateOPTEditClient Returns an OPT-175B client capable of making code edits.
func CreateOPTEditClient() ai.EditClient {
	return gptjClient{}
}
