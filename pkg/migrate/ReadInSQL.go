package migrate

import (
	"bufio"
	"io"
	"strings"
)

// ReadInSQL reads the SQL lines from the in reader and writes them to the out Buffer.
// If dropComments is true, then drops all line comments.
// If dropBlankLines is true, then drops all blank lines.
func ReadInSQL(in io.Reader, out *Buffer, dropComments bool, dropBlankLines bool, dropSearchPath bool) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "--") {
			if dropComments {
				continue
			}
		} else if strings.Contains(line, "pg_catalog.set_config('search_path'") {
			if dropSearchPath {
				continue
			}
		} else if len(strings.TrimSpace(line)) == 0 {
			if dropBlankLines {
				continue
			}
		}
		out.WriteString(line)
		out.WriteByte('\n')
	}
	out.Close()
}
