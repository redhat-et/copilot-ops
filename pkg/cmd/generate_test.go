package cmd_test

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Generate command", func() {
	// var c *cobra.Command
	// var ts *httptest.Server

	// BeforeEach(func() {
	// 	c = cmd.NewGenerateCmd()
	// })

	// When("server is created", func() {
	// 	BeforeEach(func() {
	// 		ts = OpenAITestServer()

	// 		Expect(c).NotTo(BeNil())
	// 		err := c.Flags().Set(cmd.FlagNTokensFull, "1")
	// 		Expect(err).To(BeNil())
	// 	})

	// 	JustBeforeEach(func() {
	// 		ts.Start()
	// 		log.Println("using server URL: ", ts.URL)
	// 		err := c.Flags().Set(cmd.FlagOpenAIURLFull, ts.URL+gpt3.OpenAIEndpointV1)
	// 		Expect(err).To(BeNil())
	// 		err = c.Flags().Set(cmd.FlagAIBackendFull, string(ai.GPT3))
	// 		Expect(err).To(BeNil())
	// 	})
	// 	AfterEach(func() {
	// 		ts.Close()
	// 	})

	// 	It("executes properly", func() {
	// 		err := cmd.RunGenerate(c, []string{})
	// 		// use the minimum amount of tokens from OpenAI
	// 		Expect(err).To(BeNil())
	// 	})
	// 	// TODO: add more tests for expected success
	// })

	// When("OpenAI server is down", func() {
	// 	// BeforeEach(func() {
	// 	// 	// set a port that isn't taken
	// 	// 	err := c.Flags().Set(cmd.FlagOpenAIURLFull, "http://localhost:23423")
	// 	// 	Expect(err).To(BeNil())
	// 	// })
	// 	// It("fails", func() {
	// 	// 	err := cmd.RunGenerate(c, []string{})
	// 	// 	Expect(err).To(HaveOccurred())
	// 	// })
	// 	// // TODO: add more cases that should fail
	// })
})
