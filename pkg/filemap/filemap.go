// filemap Defines the Filemap type and its methods.
package filemap

// File represents a file which was referenced in the issue to be updated.
type File struct {
	Path           string `json:"path"`
	Content        string `json:"content"`
	UpdatedContent string `json:"updatedContent"`
	Tag            string `json:"tag"`
}

// Filemap represents a mapping of files in a directory by their tagnames.
type Filemap struct {
	Files   map[string]File `json:"files"`
	Dirpath string          `json:"dirpath"`
}

// PopulateFromIssue Populates populates the filemap with the referenced files from the issue.
func (fm *Filemap) PopulateFromIssue(issue string) error {
	return nil
}

// WriteUpdatesToFiles Writes the updated contents of each file to the directory.
func (fm *Filemap) WriteUpdatesToFiles() error {
	return nil
}

// UpdateFilemapContent Decodes the given content and updates the filemap with the decoded content.
// If new files exist within the content, a best-guess effort will be made to determine the name and pathing.
func (fm *Filemap) UpdateFilemapFromContent(content string) error {
	return nil
}

// EncodeToInputText Encodes the filemap into a string which can be used as input to the OpenAI CLI.
// If there was some issue or problem encoding the filemap, an error will be returned.
func (fm *Filemap) EncodeToInputText() (string, error) {
	return "", nil
}
