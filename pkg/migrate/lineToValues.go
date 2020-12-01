package migrate

import (
	"strings"
)

func lineToValues(line string) []interface{} {
	parts := strings.Split(line, "\t")
	values := make([]interface{}, 0, len(parts))
	for _, s := range parts {
		if s == "\\N" {
			values = append(values, nil)
		} else {
			values = append(values, s)
		}
	}
	return values
}
