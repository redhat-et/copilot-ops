// filemap Defines the Filemap type and its methods.
package filemap

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/redhat-et/copilot-ops/pkg/openai"
)

// FILE Constants define the constant values that are used for parsing files.
const (
	// FILE_DELIMETER Is the string used to separate files when encoding/decoding.
	FILE_DELIMETER = "==="
	// FILE_TAG_PREFIX Is a string that indicates that the following string is the file's tag.
	FILE_TAG_PREFIX = "@"
)

// File represents a file which was referenced in the issue to be updated.
type File struct {
	// Path is the path to the file.
	Path string `json:"path"`
	// Content is the content of the file.
	Content string `json:"content"`
	// Tag is the tagname of the file.
	Tag string `json:"tag"`
	// Name is the name of the file.
	Name string `json:"name"`
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
		l := len(f.Content)
		if l > 30 {
			l = 30
		}
		short := strings.ReplaceAll(f.Content[:l], "\n", " ")
		log.Printf(" - tag: %-10q: path: %-20q [%s ...] len %d\n", f.Tag, f.Path, short, len(f.Content))
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
	for _, file := range fm.Files {
		// add extension if necessary, assume this is YAML for the time being
		// HACK: classify the relevant extension (e.g. .yaml, .yml, .json)
		// fileName := file.Tag
		// if len(strings.Split(file.Tag, ".")) == 1 {
		// 	fileName += ".yaml"
		// }
		log.Printf("path: %q, tag: %q\n", file.Path, file.Tag)
		// locate the base directory of filePath
		dirPath := filepath.Dir(file.Path)
		// create the directory if it does not exist
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				return err
			}
		}

		// write the file at the given path with read write permissions for user, read-only for others
		log.Printf("writing to file %q\n", file.Path)
		f, err := os.OpenFile(file.Path, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		// write the file
		_, err = f.WriteString(file.Content)
		if err != nil {
			return err
		}

	}
	return nil
}

// EncodeToInputText Encodes the filemap into a string which can be used as input to the OpenAI CLI.
// If there was some issue or problem encoding the filemap, an error will be returned.
func (fm *Filemap) EncodeToInputText() (string, error) {
	/*
		This function will encode the file contents as a string,
		with each file prepended by a hashtag, followed by its tagname.

		Example:
			# file1.yaml
			kind: ConfigMap
			metadata:
				name: file1
				namespace: default
			{FileDelimeter}
			# file2.yaml
			kind: ConfigMap
			metadata:
				name: file2
				namespace: default
			{FileDelimeter}
			# file3.yaml
			kind: Pod
			metadata:
				name: my-sql-pod
				namespace: default
	*/
	var input string = ""
	var i int = 0
	// join the files together along with their tag
	for tagname, file := range fm.Files {
		input += fmt.Sprintf("# %s%s\n%s\n", FILE_TAG_PREFIX, tagname, file.Content)
		// insert a delimeter between each file, but not after the last file
		if 1 < len(fm.Files) && i < len(fm.Files)-1 {
			input += fmt.Sprintf("%s\n", FILE_DELIMETER)
		}
		i++
	}
	return input, nil
}

// EncodeToInputTextFullPaths Encodes the filemap into a string using each file's full path as its tagname.
func (fm *Filemap) EncodeToInputTextFullPaths(outputType string) (string, error) {
	var input string = ""
	var i int = 0
	var GenFiles []openai.FileOutput

	// join the files together along with their tag
	for _, file := range fm.Files {

		sinFile := openai.FileOutput{
			Name:     file.Name,
			Path:     file.Path,
			Contents: file.Content,
		}

		GenFiles = append(GenFiles, sinFile)

		input += fmt.Sprintf("# %s%s\n%s\n", FILE_TAG_PREFIX, sinFile.Path, sinFile.Contents)
		// insert a delimeter between each file, but not after the last file
		if 1 < len(fm.Files) && i < len(fm.Files)-1 {
			input += fmt.Sprintf("%s\n", FILE_DELIMETER)
		}
		i++
	}

	GeneratedFiles := openai.GeneratedFilesOutput{
		GeneratedFiles: GenFiles,
	}

	genFilesOut, err := json.MarshalIndent(GeneratedFiles, "", "    ")
	if err != nil {
		return "", err
	}

	if outputType == "json" {
		return string(genFilesOut), nil
	} else if outputType == "plain" {
		return input, nil
	} else {
		err = errors.New("invalid output type")
		return "", err
	}

	//return string(genFilesOut), nil
}

// ExtractTagName Extracts the tagname from the given content, providing its line number in the content, or an error if it doesn't exist.
func ExtractTagName(content string) (string, int32, error) {
	// The tagname would be on a line in the format of: "# {FILE_TAG_PREFIX}tagname\n"
	// We can split the line by the '#' character and then trim the leading and trailing whitespace.
	lines := strings.Split(content, "\n")

	// search content for regex of the following pattern: /#\s*\{FILE_TAG_PREFIX}(.+)/g
	// if found, return the tagname
	// if not found, return an error
	pattern := fmt.Sprintf(`#\s*\%s(.+)`, FILE_TAG_PREFIX)
	regexPattern, err := regexp.Compile(pattern)
	if err != nil {
		return "", 0, err
	}

	for i, line := range lines {
		// find the first line that matches the pattern
		if match := regexPattern.FindStringSubmatch(line); match != nil {
			return strings.TrimSpace(match[1]), int32(i), nil
		}
	}
	return "", -1, fmt.Errorf("no tagname found in content")
}

// ConcatenateAfterLineNum Concatenates all of the content following the given lineNum.
// If the lineNum exceeds the number of lines in the content, an error will be returned.
func ConcatenateAfterLineNum(content string, lineNum int32) (string, error) {
	lines := strings.Split(content, "\n")
	if lineNum >= int32(len(lines)) {
		return "", fmt.Errorf("line number %d exceeds number of lines in content", lineNum)
	}
	var output string = ""
	for i := lineNum + 1; i < int32(len(lines)); i++ {
		output += lines[i] + "\n"
	}
	return output, nil
}

// AddContentByTag Adds the given content to the filemap, using the given tagname.
// If the tagname already exists, the content will be appended to the existing content.
// If the tagname does not exist, the content will be added as a new file.
func (fm *Filemap) AddContentByTag(tagname string, content string) {
	// check if the tagname already exists
	if existingFile, ok := fm.Files[tagname]; ok {
		existingFile.Content = content
		fm.Files[tagname] = existingFile
	} else {
		// if the tagname doesn't exist, add it as a new file
		fm.Files[tagname] = File{
			// TODO: infer path
			Path:    "",
			Content: content,
			Tag:     tagname,
		}
	}
}

// DecodeFromOutput Decodes the given content and updates the filemap with the decoded content.
// If new files exist within the content, a best-guess effort will be made to determine the name and pathing.
func (fm *Filemap) DecodeFromOutput(content string) error {
	// To decode from the output, we have to split the content up by the defined file delimeter.
	// Then we use RegEx to extract the tagname which we can use to look up the file and update its content.
	// If the tagname is not found, we assume that the file is new and we will create a new file with the tagname.

	// Split the content by the file delimeter
	parts := strings.Split(content, FILE_DELIMETER)
	fmt.Printf("num parts: %d\n", len(parts))
	for _, part := range parts {
		// Trim the leading and trailing whitespace
		part = strings.TrimSpace(part)
		// If the part is empty, skip it
		if part == "" {
			fmt.Printf("part is empty, skipping\n")
			continue
		}
		// The tagName should be the first line
		tagName, lineNum, err := ExtractTagName(part)
		if err != nil {
			return err
		}

		// grab the content following the lineNum
		content, err := ConcatenateAfterLineNum(part, lineNum)
		if err != nil {
			return err
		}
		// add to the filemap
		fm.AddContentByTag(tagName, content)
	}
	return nil
}
