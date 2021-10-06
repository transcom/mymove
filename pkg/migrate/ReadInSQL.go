package migrate

import (
	"strings"
	"unicode"
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

	// When loading data with `COPY ... FROM stdin;` we can have trailing tabs that are significant
	// if a record has an empty string value for the last column.
	// We could also have preceding whitespace that is significant if the first column value is empty,
	// but that has not happened in practice yet.
	return strings.TrimLeftFunc(line, unicode.IsSpace)
}
