package cmd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmd "github.com/redhat-et/copilot-ops/pkg/cmd"
)

var _ = Describe("Config", func() {
	var config *cmd.Config
	When("config is being loaded", func() {
		BeforeEach(func() {
			// generate an empty config
			config = &cmd.Config{}
		})
		When("filesets are provided", func() {
			BeforeEach(func() {
				config.Filesets = []cmd.ConfigFilesets{
					{
						Name:  "test",
						Files: []string{"test.txt"},
					},
				}
			})

			It("finds the correct filesets", func() {
				// config should find a fileset named "test"
				Expect(config.FindFileset("test")).NotTo(BeNil())
				// config should not find a fileset named "test2"
				Expect(config.FindFileset("test2")).To(BeNil())
			})

			It("is case sensitive", func() {
				// config should not find a fileset named "test"
				Expect(config.FindFileset("test")).NotTo(BeNil())
				// config should find a fileset named "TEST"
				Expect(config.FindFileset("TEST")).To(BeNil())
			})
		})
	})
})
