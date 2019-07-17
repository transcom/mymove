package migrate

import (
	"bufio"
	"bytes"
	"io"
)

var (
	lineCommentBytes   = []byte("--")
	startCopyFromStdin = []byte(" FROM stdin;")
	endCopyFromStdin   = []byte("\\.")
)

// ReadInSQL reads the SQL lines from the in reader and writes them to the out Buffer.
// If dropComments is true, then drops all line comments.
// If dropBlankLines is true, then drops all blank lines.
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

			// closing parentheses aren't strictly required and wouldn't hit on this line if there was a space after
			// search_path, so it has intentionally been left off
			if bytes.Contains(line, []byte("pg_catalog.set_config('search_path'")) {
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
