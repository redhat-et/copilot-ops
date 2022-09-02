package gpt3_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gogpt "github.com/sashabaranov/go-gpt3"

	"github.com/redhat-et/copilot-ops/pkg/ai"
	"github.com/redhat-et/copilot-ops/pkg/ai/gpt3"
)

var _ = Describe("Gpt3 Client", func() {
	var gpt3Client ai.Client

	BeforeEach(func() {
		// create client
		gpt3Client = gpt3.CreateGPT3Client("abc", nil, "")
	})

	It("doesn't generate with an empty URL", func() {
		responses, err := gpt3Client.Generate(
			gpt3.GenerateParams{
				Params: gogpt.CompletionRequest{
					Model:  "davinci",
					Prompt: "test",
				},
			},
		)
		Expect(err).To(HaveOccurred())
		Expect(responses).To(BeEmpty())
	})
})
