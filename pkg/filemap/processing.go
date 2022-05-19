// processing.go aims to provide utilities for us to recursively parse
// a given YAML file and extract the contents for any arbitrary YAML.
package filemap

import (
	"fmt"

	yamlpkg "gopkg.in/yaml.v2"
)

func ExtractYAMLKeys(yaml string) ([]string, error) {
	// parse the YAML
	var yamlMap map[string]interface{}
	// err := yaml.Unmarshal([]byte(yaml), &yamlMap)
	// if err != nil {
	// return nil, err
	// }

	testYAML := `
---
foo:
	bar:
baz:
	- 1
	- 2
	- 3
newField:
	subField:
		gary:
			age: 23
			name: gary
			children:
				- name: mary
					age: 3
				- name: larry
					age: 4
`

	yamlpkg.Unmarshal([]byte(testYAML), &yamlMap)

	fmt.Printf("%v\n", yamlMap)

	return nil, nil
}
