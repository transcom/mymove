package migrate

import (
	"strings"
)

func parseColumns(str string) []string {
	tokens := strings.Split(str, ",")
	columns := make([]string, 0, len(tokens))
	for _, c := range tokens {
		columns = append(columns, strings.TrimSpace(c))
	}
	return columns
}
