package filemap_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFilemap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filemap Suite")
}
