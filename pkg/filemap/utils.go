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
