package migrate

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var cleanMigrationPattern = regexp.MustCompile("^(\\s*)(COPY)(\\s+)([a-zA-Z._]+)(\\s*)\\((.+)\\)(\\s+)(FROM)(\\s+)(stdin)(;)(\\s*)$")

// CleanMigraton cleans migrations so they can be used in a go connection directly.
func CleanMigraton(r io.Reader, w chan string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "--") {
			continue
		} else if strings.Contains(line, "pg_catalog.set_config('search_path'") {
			continue
		} else if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		match := cleanMigrationPattern.FindStringSubmatch(line)
		if match == nil {
			w <- line
		} else {
			// 0 : Full Line
			// 1 : Whitespace Prefix
			// 2 : COPY
			// 3 : Whitespace
			// 4 : table name
			// 5 : whitespace
			// 6 : list of columns
			// 7 : whitespace
			// 8 : FROM
			// 9 : whitespace
			// 10 : stdin
			// 11 : ;
			// 12 : whitespace
			prefix := fmt.Sprintf("INSERT INTO %s (%s) VALUES\n", match[4], match[6]) // #nosec
			i := 0
			for scanner.Scan() {
				in := scanner.Text()
				if in == "\\." {
					break
				} else {
					var b strings.Builder
					b.WriteString(prefix)
					b.WriteString(" (")
					for j, s := range strings.Split(in, "\t") {
						if j != 0 {
							b.WriteString(", ")
						}
						if s == "\\N" {
							b.WriteString("NULL")
						} else if _, err := strconv.Atoi(s); err == nil {
							b.WriteString(s) // if integer don't wrap in quotes.
						} else if _, err := strconv.ParseFloat(s, 64); err == nil {
							b.WriteString(s) // if float don't wrap in quotes.
						} else {
							b.WriteString("'")
							b.WriteString(s)
							b.WriteString("'")
						}
					}
					b.WriteString(");")
					w <- b.String()
					i++
				}
			}
		}
	}
	close(w)
}
