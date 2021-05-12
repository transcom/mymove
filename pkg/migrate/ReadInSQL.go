package migrate

import (
	"strings"
)

var (
	// SQL comment designator
	sqlComment = "--"

	// closing parentheses aren't strictly required and wouldn't hit on this line if there was a space after
	// search_path, so it has intentionally been left off
	searchPath = []byte("pg_catalog.set_config('search_path'")
)

// ReadInSQLLine reads the SQL line from a string and returns the line as modified by the configuration
// If dropComments is true, then drops all line comments.
// If dropSearchPath is true, then drops all search paths.
func ReadInSQLLine(line string, dropComments bool, dropSearchPath bool) string {

	if dropComments {
		if idx := strings.Index(line, sqlComment); idx != -1 {
			line = line[0:idx]
		}
	}

	if dropSearchPath && strings.Contains(line, string(searchPath)) {
		return ""
	}

	return line
}
