package filemap_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/redhat-et/copilot-ops/pkg/filemap"
)

var _ = Describe("Filemap", func() {
	var filemap *Filemap
	BeforeEach(func() {
		// create the filemap object
		filemap = NewFilemap()
		Expect(filemap).NotTo(BeNil())
	})

	When("multiple files are loaded", func() {
		const (
			content1 = `---
kind: FortniteVod
`
			content2 = `---
kind: CatVideos
`
		)

		BeforeEach(func() {
			filemap.Files = map[string]File{
				"fortnite_vods": {
					Path:    "./testdata/fortnite_vods",
					Content: content1,
				},
				"cat_videos": {
					Path:    "./testdata/cat_videos",
					Content: content2,
				},
			}
		})

		It("contains the encoded files", func() {
			// make sure both files are contained in the filemap
			encoding := filemap.EncodeToInputText()
			Expect(encoding).To(ContainSubstring("kind: FortniteVod"))
			Expect(encoding).To(ContainSubstring("kind: CatVideos"))

			// files should be delimited by a delimiting string
			Expect(encoding).To(ContainSubstring(FileDelimeter))
		})

		It("adds new files", func() {
			// make sure the new content is added and is encoded
			filemap.AddContentByTag("new_tag", "new_content")
			encoding := filemap.EncodeToInputText()
			Expect(encoding).To(ContainSubstring("new_content"))
			Expect(encoding).To(ContainSubstring(FileDelimeter))
		})

		It("updates existing files by their tagname", func() {
			// update the fortnite vods file with content
			filemap.AddContentByTag("fortnite_vods", "new-fortnite-content")
			Expect(filemap.Files["fortnite_vods"].Content).To(ContainSubstring("new-fortnite-content"))

			// make sure that content with new tags are simply appended to the filemap
			filemap.AddContentByTag("new_tag", "new_content")
			Expect(filemap.Files["new_tag"].Content).To(ContainSubstring("new_content"))
		})

		When("the filemap is encoded to output text", func() {
			It("encodes files using their full paths", func() {
				output, err := filemap.EncodeToInputTextFullPaths(OutputPlain)
				Expect(err).NotTo(HaveOccurred())
				Expect(output).To(ContainSubstring("# @./testdata/fortnite_vods"))

				output, err = filemap.EncodeToInputTextFullPaths(OutputJSON)
				Expect(err).NotTo(HaveOccurred())
				Expect(output).To(ContainSubstring("\"path\": \"./testdata/fortnite_vods\""))
				Expect(output).To(ContainSubstring("\"path\": \"./testdata/cat_videos\""))
				Expect(output).To(ContainSubstring("kind: FortniteVod"))
				Expect(output).To(ContainSubstring("kind: CatVideos"))
			})

			It("encodes all formats except unknown", func() {
				// create an anonymous struct & iterate through to make sure that only
				// known formats are encoded
				var formats = []struct {
					Format       string `json:"format"`
					ShouldEncode bool   `json:"should_encode"`
				}{
					{Format: OutputPlain, ShouldEncode: true},
					{Format: OutputJSON, ShouldEncode: true},
					{Format: "unknown format", ShouldEncode: false},
				}
				for _, format := range formats {
					_, err := filemap.EncodeToInputTextFullPaths(format.Format)
					if format.ShouldEncode {
						Expect(err).NotTo(HaveOccurred())
					} else {
						Expect(err).To(HaveOccurred())
					}
				}
			})
		})
	})

	It("concatenates after a line number", func() {
		const content = `1
2
3
4
5`
		// check that concatenate after line number 0 includes the whole file
		cat, err := ConcatenateAfterLineNum(content, -1)
		Expect(err).NotTo(HaveOccurred())
		Expect(cat).To(ContainSubstring(content))

		// after lin num should exclude it.
		cat, err = ConcatenateAfterLineNum(content, 0)
		Expect(err).NotTo(HaveOccurred())
		Expect(cat).To(ContainSubstring(content[2:]))

		// line 5 should return an error
		_, err = ConcatenateAfterLineNum(content, 5)
		Expect(err).To(HaveOccurred())
	})

	When("decoding from an output", func() {
		const responseTemplate = `
# %sfortnite-stats
---
kind: Deployment
metadata:
	name: fortnite-stats
%s
---
# %sviva-pinata-server
kind: Deployment
apiVersion: apps/v1
metadata:
	name: viva-pinata-server
`

		It("updates the filemap from the given yamls", func() {
			var response = fmt.Sprintf(responseTemplate, FileTagPrefix, FileDelimeter, FileTagPrefix)
			err := filemap.DecodeFromOutput(response)
			Expect(err).NotTo(HaveOccurred())
			Expect(filemap.Files).To(HaveLen(2))
			Expect(filemap.Files["fortnite-stats"].Content).To(ContainSubstring("kind: Deployment"))
			Expect(filemap.Files["viva-pinata-server"].Content).To(ContainSubstring("kind: Deployment"))
			// delimeter should not be included in the content
			Expect(filemap.Files["fortnite-stats"].Content).NotTo(ContainSubstring(FileDelimeter))
			Expect(filemap.Files["viva-pinata-server"].Content).NotTo(ContainSubstring(FileDelimeter))
		})

		It("doesn't decode without tagname", func() {
			var response = fmt.Sprintf(responseTemplate, "", FileDelimeter, FileTagPrefix)
			err := filemap.DecodeFromOutput(response)
			Expect(err).To(HaveOccurred())

			response = fmt.Sprintf(responseTemplate, FileTagPrefix, FileDelimeter, "")
			err = filemap.DecodeFromOutput(response)
			Expect(err).To(HaveOccurred())
		})

		It("doesn't add empty outputs", func() {
			response := fmt.Sprintf(`# %semtpy-file
%s
# %snot-empty
kind: NotEmtpy`, FileTagPrefix, FileDelimeter, FileTagPrefix)
			err := filemap.DecodeFromOutput(response)
			Expect(err).NotTo(HaveOccurred())
			Expect(filemap.Files).To(HaveLen(1))
			// empty-file should not be in filemap
			_, ok := filemap.Files["empty-file"]
			Expect(ok).To(BeFalse())
			// not-empty should be in filemap
			listing, ok := filemap.Files["not-empty"]
			Expect(ok).To(BeTrue())
			Expect(listing.Content).To(ContainSubstring("kind: NotEmtpy"))
		})
	})
})
