package gpt3_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGpt3(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gpt3 Suite")
}
