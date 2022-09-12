package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/redhat-et/copilot-ops/pkg/cmd/config"
)

var _ = Describe("Config", func() {
	var conf *config.Config
	When("config is being loaded", func() {
		BeforeEach(func() {
			// generate an empty config
			conf = &config.Config{}
		})
		When("filesets are provided", func() {
			BeforeEach(func() {
				conf.Filesets = []config.Filesets{
					{
						Name:  "test",
						Files: []string{"test.txt"},
					},
				}
			})

			It("finds the correct filesets", func() {
				// config should find a fileset named "test"
				Expect(conf.FindFileset("test")).NotTo(BeNil())
				// config should not find a fileset named "test2"
				Expect(conf.FindFileset("test2")).To(BeNil())
			})

			It("is case sensitive", func() {
				// config should not find a fileset named "test"
				Expect(conf.FindFileset("test")).NotTo(BeNil())
				// config should find a fileset named "TEST"
				Expect(conf.FindFileset("TEST")).To(BeNil())
			})
		})
	})
})
