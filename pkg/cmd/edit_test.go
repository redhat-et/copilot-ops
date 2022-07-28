package cmd_test

import (
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-et/copilot-ops/pkg/cmd"
	"github.com/spf13/cobra"
)

var _ = Describe("Edit", func() {
	var c *cobra.Command

	BeforeEach(func() {
		// create command
		c = cmd.NewEditCmd()
		Expect(c).NotTo(BeNil())
	})

	When("OpenAI server exists", func() {
		var ts *httptest.Server
		BeforeEach(func() {
			ts = OpenAITestServer()
		})

		JustBeforeEach(func() {
			ts.Start()
			err := c.Flags().Set(cmd.FlagOpenAIURLFull, ts.URL)
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			defer ts.Close()
		})

		It("works", func() {
			err := cmd.RunEdit(c, []string{})
			Expect(err).To(BeNil())
		})

	})
})
