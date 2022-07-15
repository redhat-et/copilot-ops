package cmd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmd "github.com/redhat-et/copilot-ops/pkg/cmd"
)

var _ = Describe("Root command", func() {
	When("root command is created", func() {
		It("contains edit and generate", func() {
			rootCmd := cmd.NewRootCmd()
			Expect(rootCmd).NotTo(BeNil())
		})
	})
})
