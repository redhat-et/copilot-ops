// cmd Describes the CLI definitions and the functions which are called when the CLI is executed.
package cmd

import fm "github.com/redhat-et/copilot-ops/pkg/filemap"

// FilemapBuilder builds  given an issue and a path to a repository, build a mapping between the filetag and the file
func FilemapBuilder(issue string, dirPath string) fm.Filemap {
	filemap := fm.Filemap{}
	return filemap
}
