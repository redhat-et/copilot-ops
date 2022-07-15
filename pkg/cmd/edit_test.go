package cmd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-et/copilot-ops/pkg/cmd"
)

var _ = Describe("Edit", func() {
	It("no-op", func() {
		Expect(cmd.BuildOpenAIClient(cmd.Config{}, 0, 0, "")).NotTo(BeNil())
	})
})
