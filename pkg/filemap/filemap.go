// filemap Defines the Filemap type and its methods.
package filemap

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// File represents a file which was referenced in the issue to be updated.
type File struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

// Filemap represents a mapping of files in a directory by their tagnames.
type Filemap struct {
	Files map[string]File `json:"files"`
}

func NewFilemap() *Filemap {
	return &Filemap{
		Files: make(map[string]File),
	}
}
func (fm *Filemap) LogDump() {
	log.Printf("filemap: len %d\n", len(fm.Files))
	for _, f := range fm.Files {
		short := strings.ReplaceAll(f.Content[0:30], "\n", " ")
		log.Printf(" - %-10s: %-20s [%s ...] len %d\n", f.Tag, f.Path, short, len(f.Content))
	}
}

// LoadFilesFromGlob reads files into the filemap from the given glob pattern.
func (fm *Filemap) LoadFile(path string) error {
	tag := filepath.Base(path)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if _, ok := fm.Files[tag]; ok {
		tag = fmt.Sprintf("%s#%d", tag, len(fm.Files))
		if _, ok := fm.Files[tag]; ok {
			return fmt.Errorf("File tag conflict %s", tag)
		}
	}
	fm.Files[tag] = File{
		Path:    path,
		Content: string(bytes),
		Tag:     tag,
	}
	return nil
}

// LoadFilesFromGlob reads files into the filemap from the given glob pattern.
func (fm *Filemap) LoadFilesFromGlob(glob string) error {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	log.Printf("LoadFilesFromGlob %q - matches %v\n", glob, matches)
	for _, match := range matches {
		err := fm.LoadFile(match)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteUpdatesToFiles Writes the updated contents of each file to the directory.
func (fm *Filemap) WriteUpdatesToFiles() error {
	return nil
}

// EncodeToInputText Encodes the filemap into a string which can be used as input to the OpenAI CLI.
// If there was some issue or problem encoding the filemap, an error will be returned.
func (fm *Filemap) EncodeToInputText() (string, error) {
	return "", nil
}

// DecodeFromOutput Decodes the given content and updates the filemap with the decoded content.
// If new files exist within the content, a best-guess effort will be made to determine the name and pathing.
func (fm *Filemap) DecodeFromOutput(content string) error {
	return nil
}
