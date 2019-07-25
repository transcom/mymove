package migrate

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

var (
	lineCommentBytes   = []byte("--")
	startCopyFromStdin = []byte(" FROM stdin;")
	endCopyFromStdin   = []byte("\\.")

	// closing parentheses aren't strictly required and wouldn't hit on this line if there was a space after
	// search_path, so it has intentionally been left off
	searchPath = []byte("pg_catalog.set_config('search_path'")
)

// ReadInSQL reads the SQL lines from the in reader and writes them to the out Buffer.
// If dropComments is true, then drops all line comments.
// If dropBlankLines is true, then drops all blank lines.
// If dropSearchPath is true, then drops all search paths.
func ReadInSQL(in io.Reader, out *Buffer, dropComments bool, dropBlankLines bool, dropSearchPath bool) {
	scanner := bufio.NewScanner(in)
	inCopyFrom := false
	for scanner.Scan() {
		line := scanner.Bytes()

		if inCopyFrom {

			if bytes.Equal(line, endCopyFromStdin) {
				inCopyFrom = false
			}

		} else {

			if dropComments {
				if idx := bytes.Index(line, lineCommentBytes); idx != -1 {
					if idx == 0 {
						continue
					} else {
						line = line[0:idx]
					}
				}
			}

			if bytes.Contains(line, searchPath) {
				if dropSearchPath {
					continue
				}
			}

			if dropBlankLines {
				if len(bytes.TrimSpace(line)) == 0 {
					continue
				}
			}

			if bytes.HasSuffix(line, startCopyFromStdin) {
				inCopyFrom = true
			}

		}

		out.WriteString(string(line))
		out.WriteByte('\n')
	}
	out.Close()
}

// ReadInSQLLine reads the SQL line from a string and returns the line as modified by the configuration
// If dropComments is true, then drops all line comments.
// If dropSearchPath is true, then drops all search paths.
func ReadInSQLLine(line string, dropComments bool, dropSearchPath bool) string {

	if dropComments {
		if idx := strings.Index(line, string(lineCommentBytes)); idx != -1 {
			line = line[0:idx]
		}
	}

	if dropSearchPath && strings.Contains(line, string(searchPath)) {
		return ""
	}

	return strings.TrimSpace(line)
}
