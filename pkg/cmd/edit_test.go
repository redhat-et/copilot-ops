package cmd_test

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Edit", func() {
	// var c *cobra.Command

	// BeforeEach(func() {
	// 	// create command
	// 	c = cmd.NewEditCmd()
	// 	Expect(c).NotTo(BeNil())
	// })

	// When("OpenAI server exists", func() {
	// 	var ts *httptest.Server
	// 	BeforeEach(func() {
	// 		ts = OpenAITestServer()
	// 	})

	// 	JustBeforeEach(func() {
	// 		ts.Start()
	// 		err := c.Flags().Set(cmd.FlagOpenAIURLFull, ts.URL+gpt3.OpenAIEndpointV1)
	// 		Expect(err).To(BeNil())
	// 	})

	// 	AfterEach(func() {
	// 		defer ts.Close()
	// 	})

	// 	It("works", func() {
	// 		log.Printf("requesting the following url: %q\n", ts.URL)
	// 		err := cmd.RunEdit(c, []string{})
	// 		Expect(err).To(BeNil())
	// 	})

	// })
})
