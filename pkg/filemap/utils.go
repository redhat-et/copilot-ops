package filemap

import (
	"fmt"
	"regexp"
	"strings"
)

// extractTagName Extracts the tagname from the given content,
// providing its line number in the content, or an error if it doesn't exist.
func extractTagName(content string) (string, int32, error) {
	// The tagname would be on a line in the format of: "# {FILE_TAG_PREFIX}tagname\n"
	// We can split the line by the '#' character and then trim the leading and trailing whitespace.
	lines := strings.Split(content, "\n")

	// search content for regex of the following pattern: /#\s*\{FILE_TAG_PREFIX}(.+)/g
	// if found, return the tagname
	// if not found, return an error
	pattern := fmt.Sprintf(`#\s*\%s(.+)`, FileTagPrefix)
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

// // outputToFile Converts a given output into a file object.
// func outputToFile(output string) (File, error) {
// 	// Trim the leading and trailing whitespace
// 	part := strings.TrimSpace(output)
// 	// If the part is empty, skip it
// 	if part == "" {
// 		return File{}, fmt.Errorf("part is empty")
// 	}
// 	// The tagName should be the first line
// 	tagName, lineNum, err := extractTagName(part)
// 	if err != nil {
// 		return File{}, err
// 	}

// 	// grab the content following the lineNum
// 	concatenatedContent, err := ConcatenateAfterLineNum(part, lineNum)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// ignore empty files
// 	if strings.TrimSpace(concatenatedContent) == "" {
// 		continue
// 	}

// }

// // OutputToFile Converts a given string consisting of YAML files and returns a list of files.
// func OutputToFile(output string) ([]File, error) {
// 	// To decode from the output, we have to split the content up by the defined file delimeter.
// 	// Then we use RegEx to extract the tagname which we can use to look up the file and update its content.
// 	// If the tagname is not found, we assume that the file is new and we will create a new file with the tagname.
// 	files := make([]File, 0)

// 	// Split the content by the file delimeter
// 	parts := strings.Split(output, FileDelimeter)
// 	for _, part := range parts {

// 		// add to the filemap
// 		file := File{
// 			Name:    tagName,
// 			Path:    "generated-by-copilot-ops/",
// 			Content: concatenatedContent,
// 		}
// 		files = append(files, file)
// 	}

// 	return files, nil
// }
