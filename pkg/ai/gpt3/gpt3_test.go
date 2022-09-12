package gpt3_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-et/copilot-ops/pkg/ai"
	"github.com/redhat-et/copilot-ops/pkg/ai/gpt3"
)

var _ = Describe("Gpt3 Generate Client", func() {
	var gpt3Client ai.GenerateClient

	BeforeEach(func() {
		// create client
		gpt3Client = gpt3.CreateGPT3GenerateClient(
			gpt3.OpenAIConfig{
				Token:   "abc",
				OrgID:   nil,
				BaseURL: "http://example.com",
			},
			"hello world",
			256,
			1,
		)
	})

	It("doesn't generate with an empty URL", func() {
		responses, err := gpt3Client.Generate()
		Expect(err).To(HaveOccurred())
		Expect(responses).To(BeEmpty())
	})
})
